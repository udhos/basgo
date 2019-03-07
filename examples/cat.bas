100 c$=environ$("CAT")
110 print "will output contents of text file CAT=";c$
120 if len(c$)<1 then print "please set the env var CAT": end
130 print "output:"
140 open c$ for input as 1
150 while not eof(1)
160 input #1, line$
170 print line$
180 count=count+1
190 wend 
200 print count;" lines"
