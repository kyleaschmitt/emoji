#!/bin/bash
# "I hate you", is a perfectly reasonable response to this file.
EMOJI_SOURCE=https://unicode.org/Public/emoji/14.0/emoji-test.txt
EMOJI_DB_FILE=emoji_db.go
echo "var db = map[string][]string{" > $EMOJI_DB_FILE
curl -s $EMOJI_SOURCE | grep -E "^[^#].*qualified" | sed -e 's/^\([^;]*\);\s\([^#]*\)#\s\([^E]*\)\(E[0-9]\{1,\}\.[0-9]\{1,\}\)\s\([^\n]*\)/"\5": []string{"\4",\1","\2","\3"},/' -e's/\s\{2,\}/ /g' -e 's/ "/"/g' | sort -r  >> $EMOJI_DB_FILE
# match 4 is the qualification level. I should split this up so I can have
# multiple.
# This is fragile, and should be replaced...
# eventually
L=$(wc -l < $EMOJI_DB_FILE)
sed -i "${L}s/},/}}/" $EMOJI_DB_FILE
