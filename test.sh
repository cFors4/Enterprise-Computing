TEXT="I speak, therefore I am."
LANGUAGE="en-US"
VOICE="JennyNeural"
NAME=$LANGUAGE-$VOICE

echo "<?xml version=’1.0’?>" > text.xml
echo "<speak version=’1.0’ xml:lang=’$LANGUAGE’>" >> text.xml
echo " <voice xml:lang=’$LANGUAGE’ name=’$NAME’>" >> text.xml
echo $TEXT >> text.xml
echo " </voice>" >> text.xml
echo "</speak>" >> text.xml

go run texttospeech.go