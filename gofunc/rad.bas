110 rem Using _GOIMPORT and _GODECL to embed Go code within BASIC code
120 rem
130 _goimport("math")
140 _godecl("func degToRad(d float64) float64 {")
150 _godecl("    return d*math.Pi/180")
160 _godecl("}")
170 rem
180 rem Now using _GOFUNC to call that Go function from BASIC code
190 rem
200 d = 180
210 r = _gofunc("degToRad", d)
220 print d;"degrees in radians is";r
