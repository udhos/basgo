100 screen 1000
200 color 7,5:line (0,0)-(639,479)
210 _goproc("sleep1")
220 cls
230 color 2:for i=240 to 480: line (400,240)-(639,i): next
240 color 4:line (10,100)-(300,300),1,b,2
250 line (20,110)-(290,290),,bf,2
900 _goimport("time")
910 _godecl("func sleep1() { time.Sleep(1*time.Second) }")
920 _godecl("func sleep3() { time.Sleep(3*time.Second) }")
930 _goproc("sleep3")
940 end
