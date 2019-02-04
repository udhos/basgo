10 'List primes between two selected numbers
20 CLEAR :CLS :KEY OFF
30 PRINT "List primes between two selected numbers." :PRINT
40 INPUT "Input smaller number ";SMALL
50 INPUT "Input larger  number ";LARGE
60 IF LARGE < SMALL THEN SWAP LARGE,SMALL   'Swap if entered in wrong order
70 SMALL=INT(SMALL/2)*2+1                   'Ensure small number is odd
80 LARGE=INT(LARGE/2)*2-1                   'Ensure large number is odd
85 print "small=" small "large=" large
90 FOR K=SMALL TO LARGE STEP 2
100   X=SQR(K)                              'Only need to check to square root
110   F=0                                   'Count factors as a logic step
120   FOR J=3 TO X STEP 2
130     IF K/J-INT(K/J)=0 THEN LET F=F+1    'Totals up factors. F=0 is prime
140   NEXT J
150   IF F=0 THEN PRINT K; 'Prime has no factors
160 NEXT K
170 print
180 end
