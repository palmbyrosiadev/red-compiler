RANDINT 0 100
STORE num

UTIL INITBOOL "guessed"
PUSH "I have picked a number between 0 and 100. Can you guess it?"
PRINT

PUSH "Enter your guess: "
PRINT 

FUNC guess
    INPUT
    FLOAT
    STORE guess

    LOAD num
    LOAD guess
    LT
    STORE LT
    
    LOAD num
    LOAD guess
    GT
    STORE GT

    LOAD num
    LOAD guess
    EQ
    STORE EQ

    IF EQ UTIL PRINT "Correct!"
    IF LT UTIL PRINT "Higher!"
    IF GT UTIL PRINT "Lower!"

    LOAD EQ
    NOT
    STORE guessed
    
    UTIL PRINT ""
    PUSH "Enter your guess: "
    IF guessed PRINT
ENDFUNC

RUN guess guessed