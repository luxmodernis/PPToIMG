-- PPToIMG v1.1.0
-- Glissez un fichier .pptx ou .pdf sur l'icône pour extraire les images

property current_version : "1.1.0"
property github_repo : "luxmodernis/PPToIMG"

-- Récupère la dernière version publiée sur GitHub ("" si échec réseau)
on fetch_latest_version()
	try
		set api_url to "https://api.github.com/repos/" & github_repo & "/releases/latest"
		set latest_tag to do shell script "curl -s --max-time 5 '" & api_url & "' | python3 -c \"import sys,json; d=json.load(sys.stdin); print(d.get('tag_name','').lstrip('v'))\" 2>/dev/null || echo ''"
		return latest_tag
	on error
		return ""
	end try
end fetch_latest_version

-- Compare deux versions "x.y.z" ; renvoie true si b > a
on is_newer_version(a, b)
	set AppleScript's text item delimiters to "."
	set a_parts to text items of a
	set b_parts to text items of b
	set AppleScript's text item delimiters to ""
	repeat with i from 1 to 3
		set av to (item i of a_parts) as integer
		set bv to (item i of b_parts) as integer
		if bv > av then return true
		if bv < av then return false
	end repeat
	return false
end is_newer_version

-- Affiche la fenêtre de mise à jour avec tutoriel
on show_update_dialog(latest_tag)
	set btn to button returned of (display dialog "Une nouvelle version est disponible : v" & latest_tag & return & "Votre version actuelle : v" & current_version & return & return & "Comment mettre à jour :" & return & "1. Cliquez sur \"Télécharger\"" & return & "2. Téléchargez PPToIMG.app (ou .zip) depuis la page GitHub" & return & "3. Remplacez l'ancienne application par la nouvelle" & return & "4. Si un avertissement de sécurité apparaît : clic droit sur l'app → Ouvrir → Ouvrir" buttons {"Plus tard", "Télécharger"} default button "Télécharger" with title "PPToIMG — Mise à jour disponible")
	if btn is "Télécharger" then
		do shell script "open 'https://github.com/" & github_repo & "/releases/latest'"
	end if
end show_update_dialog

-- Traitement d'un fichier (extraction silencieuse, pas de dialog de version ici)
on process_file(file_path)
	set file_ext to do shell script "echo " & quoted form of file_path & " | sed 's/.*\\.//' | tr '[:upper:]' '[:lower:]'"
	set desktop_path to POSIX path of (path to desktop)
	set base_name to do shell script "basename " & quoted form of file_path & " | sed 's/\\.[^.]*$//'"
	set output_folder to desktop_path & base_name & "-images"

	if file_ext is "pptx" or file_ext is "zip" then
		-- Extraction PPTX (ZIP)
		try
			do shell script "mkdir -p " & quoted form of output_folder & " && unzip -j -o " & quoted form of file_path & " 'ppt/media/*' -d " & quoted form of output_folder
			set file_count to do shell script "ls " & quoted form of output_folder & " | wc -l | tr -d ' '"
			if file_count is "0" then
				do shell script "rm -rf " & quoted form of output_folder
				display dialog "Aucune image trouvée dans ce fichier." & return & return & "Vérifiez qu'il s'agit bien d'un fichier .pptx." buttons {"OK"} default button "OK" with icon caution
			else
				do shell script "open " & quoted form of output_folder
				display notification file_count & " image(s) extraite(s)" with title "PPToIMG ✓"
			end if
		on error err_msg
			display dialog "Erreur lors de l'extraction :" & return & return & err_msg buttons {"OK"} default button "OK" with icon stop
		end try

	else if file_ext is "pdf" then
		-- Extraction PDF (images embarquées via Python/pypdf)
		set python_script to "
import sys, os

output_dir = sys.argv[1]
pdf_path = sys.argv[2]

try:
    import pypdf
    reader = pypdf.PdfReader(pdf_path)
    count = 0
    for page_num, page in enumerate(reader.pages):
        for img_index, image in enumerate(page.images):
            name = image.name if image.name else f'image_p{page_num+1}_{img_index}.png'
            if not os.path.splitext(name)[1]:
                name += '.png'
            dest = os.path.join(output_dir, name)
            with open(dest, 'wb') as f:
                f.write(image.data)
            count += 1
    print(count)
except ImportError:
    print('NO_PYPDF')
except Exception as e:
    print('ERROR:' + str(e))
"
		try
			do shell script "mkdir -p " & quoted form of output_folder
			set result to do shell script "python3 -c " & quoted form of python_script & " " & quoted form of output_folder & " " & quoted form of file_path

			if result is "NO_PYPDF" then
				set install_btn to button returned of (display dialog "L'extraction de PDF nécessite la bibliothèque pypdf." & return & return & "Installer automatiquement ?" buttons {"Annuler", "Installer pypdf"} default button "Installer pypdf" with title "PPToIMG — Dépendance manquante")
				if install_btn is "Installer pypdf" then
					do shell script "python3 -m pip install pypdf --quiet"
					set result to do shell script "python3 -c " & quoted form of python_script & " " & quoted form of output_folder & " " & quoted form of file_path
				else
					do shell script "rm -rf " & quoted form of output_folder
					return
				end if
			end if

			if result starts with "ERROR:" then
				do shell script "rm -rf " & quoted form of output_folder
				display dialog "Erreur lors de l'extraction du PDF :" & return & return & (text 7 thru -1 of result) buttons {"OK"} default button "OK" with icon stop
				return
			end if

			set file_count to result as integer
			if file_count is 0 then
				do shell script "rm -rf " & quoted form of output_folder
				display dialog "Aucune image embarquée trouvée dans ce PDF." & return & return & "Le fichier ne contient peut-être que des graphiques vectoriels." buttons {"OK"} default button "OK" with icon caution
			else
				do shell script "open " & quoted form of output_folder
				display notification (file_count as string) & " image(s) extraite(s)" with title "PPToIMG ✓"
			end if
		on error err_msg
			do shell script "rm -rf " & quoted form of output_folder
			display dialog "Erreur :" & return & return & err_msg buttons {"OK"} default button "OK" with icon stop
		end try

	else
		display dialog "Fichier non reconnu : " & base_name & "." & file_ext & return & return & "Glissez un fichier .pptx ou .pdf." buttons {"OK"} default button "OK" with icon caution
	end if
end process_file

-- Drag & drop sur l'icône : extraction silencieuse, mise à jour affichée seulement si dispo
on open dropped_files
	repeat with dropped_file in dropped_files
		set file_path to POSIX path of dropped_file
		process_file(file_path)
	end repeat

	set latest_tag to fetch_latest_version()
	if latest_tag is not "" then
		if is_newer_version(current_version, latest_tag) then
			show_update_dialog(latest_tag)
		end if
	end if
end open

-- Clic sur l'icône dans le Dock / double-clic sans fichier : toujours afficher le statut
on run
	set latest_tag to fetch_latest_version()

	if latest_tag is "" then
		display dialog "PPToIMG — version " & current_version & return & return & "Impossible de vérifier les mises à jour (pas de connexion internet ?)." buttons {"OK"} default button "OK" with title "PPToIMG"
		return
	end if

	if is_newer_version(current_version, latest_tag) then
		show_update_dialog(latest_tag)
	else
		display dialog "PPToIMG — version " & current_version & return & return & "Vous utilisez la dernière version. ✓" & return & return & "Astuce : glissez un fichier .pptx ou .pdf sur l'icône de cette application pour en extraire les images." buttons {"OK"} default button "OK" with title "PPToIMG"
	end if
end run
