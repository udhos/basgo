10 x=9:dim stack(7):top=0
20 print "Calculating factorial for ", x, "!"
30 gosub 100
40 print x, "! = ", y
50 end

100 rem Factorial for x is returned in y
110 rem Input: x, stack(), top
120 rem Output: y
130 if x < 2 then y = 1: return
140 stack(top) = x : top = top + 1 : rem push x
160 x = x - 1
170 gosub 100
180 top = top - 1 : x = stack(top) : rem pop x
190 y = y * x 
200 return

