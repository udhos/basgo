10 dim t(16):size=16
20 print "creating deck with "size" cards"
30 for i=1 to size
40 t(i)=i
50 next
60 print "drawing "size" cards from deck"
70 n=size
80 j=1
90 while n>0
100 i=int(rnd * n) + 1
110 print "removed card "j" from deck, its value was: " t(i)
120 t(i)=t(n)
130 j=j+1
140 n=n-1
150 wend
160 end
