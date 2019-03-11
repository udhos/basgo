100 screen 0
110 print "printing all characters in current codepage"
120 for i=0 to 255
130 print chr$(i);
140 if i MOD 16 = 0 then print
150 next
200 print:print
210 print "hit any key"
220 print input$(1)
