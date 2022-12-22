UTIL INITBOOL "powering"
UTIL INITNUM "num"
UTIL INITNUM "power"
UTIL INITNUM "powc"

UTIL SET "num" 2
UTIL SET "power" 3
UTIL SET "powc" 2

LOAD num
STORE current

FUNC power
    LOAD current
    LOAD num
    MULT
    STORE current
    LOAD powc
    PUSH 1
    ADD
    STORE powc
    LOAD powc
    LOAD power
    GTE
    STORE powering
ENDFUNC

RUN power powering

UTIL PRINTVAR "current"