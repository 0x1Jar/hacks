let fs = require("fs");
let acorn = require("acorn");

let fn = process.argv[2];

if (fn == ""){
   console.log("usage: jsstrings <file>");
   process.exit() ;
}

fs.readFile(fn, "utf8", function(err, data) {
    if (err) {
        console.error("Error reading file:", err.message);
        process.exit(1);
    }
    for (let token of acorn.tokenizer(data,{
        onComment: function(block, text, start, end){
            console.log(text.trim()); // Trimmed the comment text
        }
    })) {
        if (token.type == acorn.tokTypes.string){
            console.log(token.value);
        }
    }
});
