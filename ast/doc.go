package ast

/**
思考：let x = 5; 应该如何构建 AST？
 		*ast.Program
   		  Statements
			 ↓
	  *ast.LetStatement   let
	  ---   Name
      |     Value -----
      ↓               ↓
*ast.Identifier   *ast.expression
	  x               5
*/
