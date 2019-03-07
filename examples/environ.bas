100 print "display environment table:"
110 n=1
120 e$=environ$(n)
130 if len(e$) < 1 then end
140 print "[";n;"]: ";e$
150 n=n+1
160 goto 120
