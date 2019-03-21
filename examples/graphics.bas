100 key off:cls:screen 9
200 color 7,5 : rem next CLS will clear bg color to 5
205 line (0,0)-(639,479)
206 line (50,0)-(0,50)
210 rem _goproc("sleep1")
215 print input$(1)
220 cls
230 color 2:for i=50 to 300: line (100,50)-(319,i): next
240 color 4
250 line (10,100)-(40,130),1,b
260 line (15,105)-(35,125),,bf
270 line (80,130)-(50,100),1,b
280 line (75,125)-(55,105),,bf
290 line (40,140)-(10,170),1,b
300 line (15,165)-(35,145),,bf
310 line (50,170)-(80,140),1,b
320 line (55,165)-(75,145),,bf
900 rem _goimport("time")
910 rem _godecl("func sleep1() { time.Sleep(1*time.Second) }")
920 rem _godecl("func sleep3() { time.Sleep(3*time.Second) }")
930 rem _goproc("sleep3")
935 print input$(1)
940 end
