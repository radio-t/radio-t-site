#!/bin/sh

# @raycast.schemaVersion 1
# @raycast.title radio-t process and publish mp3
# @raycast.mode silent
#
# Optional parameters:
# @raycast.icon 🎙
# @raycast.packageName Podcasts

cat <<'EOF' > /tmp/process_and_publish_mp3.sh
set -e

# Get the currently selected file in Finder
selectedFile=$(osascript -e 'tell application "Finder"' -e 'set selectedItems to selection' -e 'if (count of selectedItems) > 0 then' -e 'set selectedItem to item 1 of selectedItems' -e 'POSIX path of (selectedItem as alias)' -e 'else' -e '""' -e 'end if' -e 'end tell')

# Extract episode number from the full path
filename=$(basename "${selectedFile}")
episodeNumber=$(echo "${filename}" | sed -E 's/rt_podcast([0-9]+)\.mp3/\1/')
# Check if episode number is valid (purely numeric and not empty)
if ! [[ $episodeNumber =~ ^[0-9]+$ ]]; then
    say -v Alex "processing failed, wrong file name"
    echo "Error: Filename does not match expected pattern 'rt_podcastNNN.mp3'"
    exit 1
fi

echo "selected file: ${selectedFile}, episode number: ${episodeNumber}"

say -v Alex "starting processing for episode $episodeNumber"
cd ~/dev.umputun/radio-t.com/radio-t-site/publisher
make proc-mp3 FILE=${selectedFile}
echo "done"

osascript -e "display notification \"processing completed for $episodeNumber\" with title \"RADIO-T\""
say "completed processing for episode $episodeNumber"
EOF

# Make the script executable
chmod +x /tmp/process_and_publish_mp3.sh

# Use AppleScript to open a new Terminal window and execute the script
osascript <<END
tell application "Terminal"
    do script "/bin/sh /tmp/process_and_publish_mp3.sh"
    activate
end tell
END
