5 screen 0
10 FOR A=1 TO 10
20 T1=TIMER
30 FOR YD=1 TO 23:FOR XD=1 TO 80:LOCATE YD,XD:?CHR$(asc("0")+A);:NEXT:NEXT
40 T2=TIMER
50 T(A)=T2-T1
60 NEXT:?"Average="(T(1)+T(2)+T(3)+T(4)+T(5)+T(6)+T(7)+T(8)+T(9)+T(10))/10
1000 print "hit q to exit"
1010 i$=input$(1)
1020 while i$<>"q"
1030 print "key:";i$;" asc:";asc(i$)
1040 print "hit q to exit"
1050 i$=input$(1)
1060 wend
