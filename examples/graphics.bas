100 screen 1000
200 line (1,1)-(640,480)
900 _goimport("time")
910 _godecl("func sleep() { time.Sleep(3*time.Second) }")
920 _goproc("sleep")
930 end
