10 FOR A=1 TO 10
20 T1=TIMER
30 FOR X=1 TO 100000:NEXT
40 T2=TIMER
50 T(A)=T2-T1:?T(A)
60 NEXT:?"Average="(T(1)+T(2)+T(3)+T(4)+T(5)+T(6)+T(7)+T(8)+T(9)+T(10))/10
