01 REM https://github.com/skx/gobasic/blob/master/examples/99-game.bas

10 b=7 : rem b = RND * 100
20 LET count=1
30 PRINT "I have picked a random number, please guess it!!" 
40 PRINT "Enter your choice: ";
50 INPUT a
60 PRINT
70 IF b = a THEN 2000 ELSE PRINT "You were wrong: ";
80 IF a < b THEN PRINT "too low"
90 IF a > b THEN PRINT "too high"
100 LET count = count + 1
110 GOTO 40

2000 PRINT "You guessed my number!"
2010 PRINT "You took ", count, " attempts"
2020 END
