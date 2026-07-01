package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/sys/windows/registry"
)

const VERSION = "1.2.1"
const GITHUB_REPO = "luxmodernis/PPToIMG"

func main() {
	// Auto-enregistrement (idempotent) dans le menu clic droit de l'Explorateur.
	// Ne nécessite pas de droits admin (clé HKEY_CURRENT_USER).
	registerContextMenu()

	args := os.Args[1:]

	if len(args) == 0 {
		// Double-clic sans fichier : équivalent du clic Dock sur Mac -> toujours afficher le statut
		showStatusAndMaybeUpdate()

		// Puis on enchaîne sur le sélecteur de fichier (l'action principale sur PC)
		file := openFilePicker()
		if file != "" {
			processFile(file)
		}
	} else {
		// Fichier(s) glissé(s) sur l'icône : extraction silencieuse
		for _, f := range args {
			processFile(strings.TrimSpace(f))
		}
		// Mise à jour affichée seulement si une nouvelle version existe
		latest := fetchLatestVersion()
		if semverGreater(latest, VERSION) {
			showUpdateDialog(latest)
		}
	}
}

func processFile(filePath string) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pptx", ".zip":
		extractPPTX(filePath)
	case ".pdf":
		extractPDF(filePath)
	default:
		showDialog(
			"Fichier non reconnu :\n"+filepath.Base(filePath)+"\n\nVeuillez sélectionner un fichier .pptx ou .pdf.",
			"PPToIMG",
		)
	}
}

// ── Statut / mise à jour (clic simple sans fichier) ────────────────────────

func showStatusAndMaybeUpdate() {
	latest := fetchLatestVersion()

	if latest == "" {
		showDialog(
			fmt.Sprintf("PPToIMG — version %s\n\nImpossible de vérifier les mises à jour (pas de connexion internet ?).", VERSION),
			"PPToIMG",
		)
		return
	}

	if semverGreater(latest, VERSION) {
		showUpdateDialog(latest)
	} else {
		showDialog(
			fmt.Sprintf("PPToIMG — version %s\n\nVous utilisez la dernière version. ✓\n\nSélectionnez un fichier .pptx ou .pdf dans la fenêtre suivante pour en extraire les images.", VERSION),
			"PPToIMG",
		)
	}
}

// ── Extraction PPTX ────────────────────────────────────────────────────────

func extractPPTX(filePath string) {
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	outputFolder := filepath.Join(getDesktopPath(), baseName+"-images")

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		showDialog("Impossible de créer le dossier :\n"+err.Error(), "PPToIMG - Erreur")
		return
	}

	r, err := zip.OpenReader(filePath)
	if err != nil {
		os.RemoveAll(outputFolder)
		showDialog("Impossible de lire le fichier .pptx :\n"+err.Error(), "PPToIMG - Erreur")
		return
	}
	defer r.Close()

	count := 0
	for _, f := range r.File {
		name := f.Name
		if !strings.HasPrefix(name, "ppt/media/") || strings.HasSuffix(name, "/") {
			continue
		}
		destPath := filepath.Join(outputFolder, filepath.Base(name))
		if err := extractZipEntry(f, destPath); err == nil {
			count++
		}
	}

	if count == 0 {
		os.RemoveAll(outputFolder)
		showDialog("Aucune image trouvée dans ce fichier.\n\nVérifiez qu'il s'agit bien d'un .pptx contenant des images.", "PPToIMG")
		return
	}

	exec.Command("explorer.exe", outputFolder).Start()
	showToast("PPToIMG ✓", fmt.Sprintf("%d image(s) extraite(s)  →  %s-images", count, baseName))
}

func extractZipEntry(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

// ── Extraction PDF ─────────────────────────────────────────────────────────

func extractPDF(filePath string) {
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	outputFolder := filepath.Join(getDesktopPath(), baseName+"-images")

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		showDialog("Impossible de créer le dossier :\n"+err.Error(), "PPToIMG - Erreur")
		return
	}

	if err := api.ExtractImagesFile(filePath, outputFolder, nil, nil); err != nil {
		os.RemoveAll(outputFolder)
		showDialog("Erreur lors de l'extraction du PDF :\n\n"+err.Error(), "PPToIMG - Erreur")
		return
	}

	entries, err := os.ReadDir(outputFolder)
	if err != nil || len(entries) == 0 {
		os.RemoveAll(outputFolder)
		showDialog("Aucune image trouvée dans ce PDF.\n\nLe fichier ne contient peut-être que des graphiques vectoriels.", "PPToIMG")
		return
	}

	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}

	exec.Command("explorer.exe", outputFolder).Start()
	showToast("PPToIMG ✓", fmt.Sprintf("%d image(s) extraite(s)  →  %s-images", count, baseName))
}

// ── Auto-update ────────────────────────────────────────────────────────────

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func fetchLatestVersion() string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/" + GITHUB_REPO + "/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return ""
	}
	return strings.TrimPrefix(rel.TagName, "v")
}

func semverGreater(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	ap := strings.Split(a, ".")
	bp := strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		var av, bv int
		if i < len(ap) {
			av, _ = strconv.Atoi(ap[i])
		}
		if i < len(bp) {
			bv, _ = strconv.Atoi(bp[i])
		}
		if av > bv {
			return true
		}
		if av < bv {
			return false
		}
	}
	return false
}

func showUpdateDialog(version string) {
	msg := fmt.Sprintf(
		"Une nouvelle version est disponible : v%s\n"+
			"Votre version actuelle : v%s\n\n"+
			"Comment mettre à jour :\n"+
			"1. Cliquez sur \"Télécharger\"\n"+
			"2. Téléchargez le nouveau PPToIMG.exe depuis la page GitHub\n"+
			"3. Remplacez l'ancien fichier .exe par le nouveau\n"+
			"4. Si Windows affiche un avertissement, cliquez sur\n"+
			"    \"Informations complémentaires\" puis \"Exécuter quand même\"",
		version, VERSION,
	)
	msgEscaped := strings.ReplaceAll(msg, "'", "''")

	script := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$form = New-Object System.Windows.Forms.Form
$form.Text = 'PPToIMG - Mise a jour disponible'
$form.Size = New-Object System.Drawing.Size(460, 300)
$form.StartPosition = 'CenterScreen'
$form.FormBorderStyle = 'FixedDialog'
$form.MaximizeBox = $false
$form.MinimizeBox = $false
$form.Font = New-Object System.Drawing.Font('Segoe UI', 9)

$label = New-Object System.Windows.Forms.Label
$label.Text = '%s'
$label.SetBounds(20, 20, 410, 190)
$label.AutoSize = $false

$btnLater = New-Object System.Windows.Forms.Button
$btnLater.Text = 'Plus tard'
$btnLater.SetBounds(210, 220, 100, 32)
$btnLater.DialogResult = [System.Windows.Forms.DialogResult]::Cancel

$btnDownload = New-Object System.Windows.Forms.Button
$btnDownload.Text = 'Telecharger'
$btnDownload.SetBounds(320, 220, 110, 32)
$btnDownload.DialogResult = [System.Windows.Forms.DialogResult]::OK

$form.Controls.Add($label)
$form.Controls.Add($btnLater)
$form.Controls.Add($btnDownload)
$form.CancelButton = $btnLater
$form.AcceptButton = $btnDownload

$result = $form.ShowDialog()
if ($result -eq [System.Windows.Forms.DialogResult]::OK) {
    Start-Process 'https://github.com/%s/releases/latest'
}
`, msgEscaped, GITHUB_REPO)
	runPS(script)
}

// ── Menu clic droit "Extraire avec PPToIMG" ────────────────────────────────

func registerContextMenu() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return
	}

	for _, ext := range []string{".pptx", ".pdf"} {
		keyPath := `Software\Classes\SystemFileAssociations\` + ext + `\shell\PPToIMG`

		k, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
		if err != nil {
			continue
		}
		k.SetStringValue("", "Extraire avec PPToIMG")
		k.SetStringValue("Icon", exePath+",0")
		k.Close()

		ck, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath+`\command`, registry.SET_VALUE)
		if err != nil {
			continue
		}
		ck.SetStringValue("", `"`+exePath+`" "%1"`)
		ck.Close()
	}
}

// ── Utilitaires Windows ────────────────────────────────────────────────────

func getDesktopPath() string {
	out, err := exec.Command("reg", "query",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`,
		"/v", "Desktop").Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "Desktop") && strings.Contains(line, "REG_SZ") {
				parts := strings.SplitN(line, "REG_SZ", 2)
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Desktop")
}

func encodePSBase64(script string) string {
	utf16 := make([]byte, len([]rune(script))*2)
	i := 0
	for _, c := range script {
		utf16[i] = byte(c)
		utf16[i+1] = byte(c >> 8)
		i += 2
	}
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result []byte
	for j := 0; j < len(utf16); j += 3 {
		rem := len(utf16) - j
		var b0, b1, b2 byte
		b0 = utf16[j]
		if rem > 1 {
			b1 = utf16[j+1]
		}
		if rem > 2 {
			b2 = utf16[j+2]
		}
		result = append(result, chars[b0>>2])
		result = append(result, chars[((b0&0x3)<<4)|(b1>>4)])
		if rem > 1 {
			result = append(result, chars[((b1&0xf)<<2)|(b2>>6)])
		} else {
			result = append(result, '=')
		}
		if rem > 2 {
			result = append(result, chars[b2&0x3f])
		} else {
			result = append(result, '=')
		}
	}
	return string(result)
}

func runPS(script string) (string, error) {
	encoded := encodePSBase64(script)
	out, err := exec.Command("powershell", "-STA", "-WindowStyle", "Hidden", "-EncodedCommand", encoded).Output()
	return strings.TrimSpace(string(out)), err
}

func openFilePicker() string {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$d = New-Object System.Windows.Forms.OpenFileDialog
$d.Title = 'Selectionner un fichier'
$d.Filter = 'Fichiers pris en charge (*.pptx, *.pdf)|*.pptx;*.pdf|PowerPoint (*.pptx)|*.pptx|PDF (*.pdf)|*.pdf|Tous les fichiers (*.*)|*.*'
if ($d.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $d.FileName }
`
	result, _ := runPS(script)
	return result
}

func showDialog(message, title string) {
	msg := strings.ReplaceAll(message, "'", "''")
	ttl := strings.ReplaceAll(title, "'", "''")
	script := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
[System.Windows.Forms.MessageBox]::Show('%s', '%s', [System.Windows.Forms.MessageBoxButtons]::OK, [System.Windows.Forms.MessageBoxIcon]::Information) | Out-Null
`, msg, ttl)
	runPS(script)
}

func showToast(title, message string) {
	ttl := strings.ReplaceAll(title, "'", "''")
	msg := strings.ReplaceAll(message, "'", "''")
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType=WindowsRuntime] | Out-Null
$t = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$x = [xml]$t.GetXml()
($x.GetElementsByTagName('text'))[0].AppendChild($x.CreateTextNode('%s')) | Out-Null
($x.GetElementsByTagName('text'))[1].AppendChild($x.CreateTextNode('%s')) | Out-Null
$s = New-Object Windows.Data.Xml.Dom.XmlDocument
$s.LoadXml($x.OuterXml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('PPToIMG').Show([Windows.UI.Notifications.ToastNotification]::new($s))
`, ttl, msg)
	runPS(script)
}
