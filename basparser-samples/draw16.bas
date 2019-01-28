10 dim t(16)
20 for i=1 to 16
30 t(i)=i
40 next
50 n=16
60 while n>1
70 i=int(rnd * n) + 1
80 print t(i)
90 t(i)=t(n)
100 n=n-1
110 wend
120 end
