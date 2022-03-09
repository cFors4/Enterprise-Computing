#!/bin/sh
TEXT="I speak, therefore I am."
LANGUAGE="en-US"
VOICE="JennyNeural"
NAME=$LANGUAGE-$VOICE

echo "<?xml version='1.0'?>" > text2.xml
echo "<speak version='1.0' xml:lang='$LANGUAGE'>" >> text2.xml
echo " <voice xml:lang='$LANGUAGE' name='$NAME'>" >> text2.xml
echo "$TEXT" >> text2.xml
echo " </voice>" >> text2.xml
echo "</speak>" >> text2.xml

go run texttospeech.go