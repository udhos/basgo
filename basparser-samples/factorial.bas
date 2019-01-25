10 print "This program calculates factorial recursively using GOSUB"
20 print "Enter number: ";
30 input x
40 print "Calculating factorial for ", x, "!"
50 gosub 100
60 print "Done: ", x, "! = ", y
70 end

100 rem Factorial for x is returned in y
110 rem Input: x
120 rem Output: y
130 if x < 2 then y = 1: return
140 x = x - 1
150 gosub 100
160 x = x + 1
170 y = y * x 
180 return

