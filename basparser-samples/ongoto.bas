10 on 1 goto 100,200,300
20 a=2:on a goto 100,200,300
30 on a+1 goto 100,300,300

100 print 100 : goto 20
200 print 200 : goto 30
300 print 300
