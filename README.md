# minil
`minil` is a minimalistic language for describing test LCI models. It can be
converted to the [openLCA JSON-LD format](https://github.com/GreenDelta/olca-schema)
using the `minil` command line tool:

```bash
$ minil [input file]
```

This will produce a zip file that can be imported into
[openLCA](http://www.openlca.org/).


## Syntax
`minil` looks like this:

```r
# a simple example
p1 -> 2.2 e1
p2 <- 0.5 p1
p2 -> 0.8 w1
w1 <- 1.2 r1
w1 <- 2.0 w1
w1 -> 3.4 e1
```

It just describes the input and output relations of flows. Flows have an
identifier like `p1` where the first character should be a letter followed
by any character except whitespaces:

```ebnf
Identifier = Letter, {Character - Whitespace}
```

`minil` knows three types of flows: products, wastes, and elementary flows. The
flow type is identified by the first letter of an identifier. Idenfifiers that
start with the letter `p` (or `P`) are recognized as product flows (e.g. `p1`).
When the first letter is a `w` (or `W`) it is recognized as a waste flow
(e.g. `w1`). All other identifiers are recognized as elementary flows
(e.g. `e1` or `r1`).

```ebnf
ProductIdentifier = "p" | "P", {Character - Whitespace}
```

```ebnf
WasteIdentifier   = "w" | "W", {Character - Whitespace}
```

Each line in a `minil` file describes an input or output relation. An input
relation has the following syntax:

```ebnf
Input = ProductIdentifier | WasteIdentifier, "<-", Number, Identifier
```

For example, `w1 <- 1.2 r1` describes an input of `1.2` units of the elementary
flow `r1` for the treatment of the waste flow `w1`. Correspondingly, an output
relation has the following syntax:

```ebnf
Output = ProductIdentifier | WasteIdentifier, "->", Number, Identifier
```

The line `p2 -> 0.8 w1` describes an output of `0.8` units of the waste flow 
`w1` related to the product flow `p2`. Finally, each line that starts with a
`#` is ignored as comment.

```ebnf
Comment = "#", {Character}
```


## Conversion
For each product and waste flow a corresponding process with the same identifier
is created. For product flows, an output and for waste flows an input is added
as reference flow to that process. By default the amount of this reference
flow is set to `1.0` but it can be set to another value via an input or
output relation, e.g. sets the amount of the reference flow to `2.0` units of
inputs of the waste flow `w1`:

```r
w1 <- 2.0 w1
```

In the resulting JSON-LD package the units of all inputs and outputs are set
to kilogram mass using the UUIDs from the openLCA reference data. The converted
example above looks like the following in openLCA:



