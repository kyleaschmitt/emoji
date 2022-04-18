# emoji
A command line emoji lookup, written in golang.

I don't like to leave the terminal unless I have to, and got tired of typing in codepoints.

# usage
simple lookup
> me@debian:~$ emoji facepalm
> 
> 🤦‍♂️

more complicated lookup
> me@debian:~$ emoji woman medium light skin vampire
>
>🧛🏼‍♀️


It doesn't always get things right, but hopefully future updates will allow for globbing or saved preferences.

# mechanism
The db.go file is just a map of official unicode descriptions, to an array of 4 strings: the emoji version, the codepoints (as ascii), status, and the glyph itself.  This file is generated by the script mk_emoji_db_go.sh, which scrapes the emoji-test.txt file on unicode.org: https://unicode.org/Public/emoji/14.0/emoji-test.txt  ... and then needs some hand editing.

The emoji program uses the FuzzySearch utility from the go-edlib library to select the most likely candidate, and prints it on standard out.

There are probably lots of bugs and edge cases that aren't handled, and a good deal of functionality is still missing.  This is a rushed first pass.

