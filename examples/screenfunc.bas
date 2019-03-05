100 screen 0
110 locate 1,1: color 1,6: print "x"
120 print
130 color 7,0
140 print "code:";screen(1,1)
150 print "char:";chr$(screen(1,1))
160 attr=screen(1,1,1)
170 print "attr:";attr
180 print "fg:";attr MOD 16
190 print "bg:";attr \ 16
200 print
210 print "hit any key";input$(1)
