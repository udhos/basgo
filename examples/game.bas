10 REM ** Number guessing game **
20 REM
30 REM Adapted from:
40 REM https://github.com/skx/gobasic/blob/master/examples/99-game.bas

100 max = 200
110 b = RND * max
120 c% = b : b = c% : rem Drop decimal portion from random number
130 LET count=1
140 PRINT "I have picked a random number between 0 and "max", please guess it!"

200 REM -- Main Game Loop --
210 PRINT "Enter your choice: ";
220 INPUT a
230 PRINT
240 IF b = a THEN 2000 ELSE PRINT "You were wrong: ";
250 IF a < b THEN PRINT "too low"
260 IF a > b THEN PRINT "too high"
270 LET count = count + 1
280 GOTO 200

2000 REM -- The End --
2010 PRINT "You guessed my number!"
2020 PRINT "You took ", count, " attempts"
2030 END
