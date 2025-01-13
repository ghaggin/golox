# golox

## Grammar
### Syntax Grammar
```
program        → declaration* EOF ;
```

#### Declarations
A program is a series of declarations, which are the statements that bind new identifiers or any of the other statement types.
```
declaration    → varDecl 
               | statement ;

varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
```

#### Statements
The remaining statement rules produce side effects, but do not introduce bindings.
```
statement      → exprStmt 
               | printStmt 
               | block ;

exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
block          → "{" declaration* "}"
```

#### Expressions
Expressions produce values.

```
expression     → assignment ;
assignment     → IDENTIFIER "=" assignment | equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
```

### Lexical Grammar
The lexical grammar is used by the scanner to group characters into tokens.

```
NUMBER         → DIGIT+ ( "." DIGIT+ )? ;
STRING         → "\"" <any char except "\"">* "\"" ;
IDENTIFIER     → ALPHA ( ALPHA | DIGIT )* ;
ALPHA          → "a" ... "z" | "A" ... "Z" | "_" ;
DIGIT          → "0" ... "9" ;
```
