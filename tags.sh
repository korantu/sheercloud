[ -f TAGS ] && rm TAGS
find . -name "*.go" | xargs etags -a -r '/^\/\/---> +\([a-zA-Z]+\)/' -r '/^type \([a-zA-Z_]+\)/' -r '/^func \([a-zA-Z_]+\)/' -r '/^func ([^)]+) \([a-zA-Z_]+\)/' 
