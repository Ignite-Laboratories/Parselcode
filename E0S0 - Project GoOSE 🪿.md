# `E0S0 - Project GoOSE 🪿`
### `Alex Petz, Ignite Laboratories, November 2025`

---

### The Go OS Environment

Since I was a child, I've wanted to make my own "operating system" - but that terminology never quite fit my vision.

I want to minimize the amount of steps it takes an engineer to _access_ hardware by defining the **language** of _memory_.

A neural operating system isn't about reinventing low level memory access and hardware control, it's about _exposing_
them to the engineer through first-order retrieval (in Adam Savage's terms).

0. Let's give the language a way to share memory amongst goroutines in a thread safe manner
1. Let's give the language a way to repeat its own actions without recompilation
2. Let's expose hardware observance from memory, rather than directly interfacing with it
3. Let's make creating ad-hoc graphical windows that engineers can draw upon a 0-configuration process

Hardware _observance_ is a critical one!  Multiplexing event-like messages is the _exact_ purpose of this architecture =)

When your program launches, it should be given a _channel_ in memory through which events directed towards 
it are received.  This channel should be a _given,_ not something that each program has to derive a process
to read from - as it currently exists.  When we put that on the engineer we wind up with libraries that aim
to tackle _too much,_ and in turn become behemoths to maintain in and of themselves.

Instead - I propose a form of "standardized memory" - populated on the program's startup and used to orient it
in quickly performing standardized operating system _environmental_ tasks like graphical window creation.

This file is a scratch pad for my steps in the process of extending the language.  This, hopefully, will
help guide Other - as well as myself - in revisiting this in the future.

This, like all of my work, is an evolving stream of consciousness =)

### Sparklets

The next major component is that we need a way to describe a long-running process that's hosted
across multiple goroutines, but still a unique subset of the larger program.  The term 'module' is
already used by Go, and 'applet' doesn't really describe it well.  Instead, my chosen term for
initiating these processes is to "Spark() them  off" - so, I'm choosing to lean into that and call
my sub-processes `Sparklets`.

### Language Extensions

Assuming that you've read the Parselcode solution, let's break down what GoOSE will imbue Go with:

0. Cursor accessors and chaining of index accessors using a semicolon
1. `swizzle(target, members...)`, `parse(code...)` (and `parsel(target, members...)` combining the two)
2. `reveal(path...)`, `recall[TOut](path...)`, `describe[TOut](value, path...)` 
3. The `ref` keyword for quickly creating an inline discardable reference
4. The `cleanup` keyword for deferring until _shutdown_
5. The `rec` package for controlling _scoped_ logging, verbosity, and silent modes
6. The `core` package for performing instance-specific operations like deferring cleanup to the time of shutdown
7. The `atlas` system for configuration
8. Letter vectors
9. The `fuzzy[TOut](any)` keyword mechanic

### Environment Setup

First, if you'd like to work with a local development build of Go it's an absolute breeze.  Your IDE can easily
point to the development folder and allow you to work against it.  That being said, if you'd like to test from
your terminal, you could add the following to your `rc` file:

    ---------------------
    ~/.zshrc or ~/.bashrc
    ---------------------

    # GoOSE 🪿
    export GOROOT_RELEASE="/usr/local/go"  # Adjust to your release Go installation
    export GOROOT_DEV="$HOME/source/ignite/goose"  # Your development Go source
    
    // Wake me up, before you...
    go-go() {
        export GOROOT="$GOROOT_RELEASE"
        export PATH="$GOROOT/bin:${PATH//$GOROOT_DEV\/bin:/}"
        export GOTOOLCHAIN=local
        echo "✓ Switched to release Go toolchain"
        go version
    }
    
    go-goose() {
        export GOROOT="$GOROOT_DEV"
        export PATH="$GOROOT/bin:${PATH//$GOROOT_RELEASE\/bin:/}"
        export GOTOOLCHAIN=local
        echo "✓ Switched to GoOSE 🪿"
        go version
    }
    
    # Default to release Go on startup
    go-go > /dev/null

This allows you to run `go-goose` to switch to GoOSE in a terminal environment, or `go-go` to switch to the release version of Go.

To set up your local development branch of Go, you can clone it locally and then run `./make.bash` from the `src` folder.  If
you need to add entries to the .gitignore without commiting them (perhaps your IDE's working folder), you can do so through
editing `.git/info/exclude` in your repo's folder.  You may also fork the repository to your own remote for isolated commits,
but that's really not necessary.

For some of the work we'll be doing, you'll need to install `stringer` - please be sure to do so from the -release- environment,
not GoOSE, as this gets added to your PATH and should be installed into the appropriate GOROOT.

    go-go
    go install golang.org/x/tools/cmd/stringer@latest 

If you get an error noting that `stringer` can't be found, please add the following to your 'rc' file

    export PATH="$HOME/go/bin:$PATH"

**KEY NOTES:**

- Build the source from the `goroot/src` directory using `./make.bash`

- Configure your IDE to switch between release and your development repository through build configurations

- Modify the repository's `.git/info/exclude` to perform `.gitignore` operations specific to your local workspace

--------

### The version number

First, to ensure we can check if we are running GoOSE or Go, we get to enhance the global Go version number
with the word `GoOSE` - luckily, that's quite easy and minimally invasive to do so!  Just add a `VERSION` file
to your goroot and include the exact name of your build you'd like.  NOTE: This does not include automatic
number incrementing or anything, but that's totally fine

    VERSION
    goose1.26_devel

### Adding Tokens

The first thing we get to do is tell Go about the new keywords.  To do so, add them to 
`src/cmd/compile/internal/syntax/tokens.go` and then regenerate the token string methods in 
`src/cmd/compile/internal/syntax/token_string.go` using `src/cmd/compile/internal/syntax/go generate` or
the inline code generation helpers the developers included.

Here are a few of the tokens I've added:

    // Cursor Accessors
	_Lcaret  // <
	_Rcaret  // >
	_Pipe    // |

    // Swizzling and Parsing
	_Swizzle // swizzle
	_Parse   // parse
	_Parsel  // parsel

This allows us to use `|absolute clamp|`, `||relative clamp||`, `<absolute flow>`, and `<<relative flow>>` - as
well as `parse`, `parsel`, and `swizzle`. The double square-bracket `[[relative panic]]` tokens are already baked 
into the language through index accessors.

In addition, I've added the following tokens:

    // Standardized memory access
	_Reveal      // reveal
	_Recall      // recall
	_Describe    // describe

    // Convenience
	_Ref         // ref

    // Deferral
	_Cleanup     // cleanup

The memory access functions act as the global way to _reveal_ (untyped), _recall_ (typed), and _describe_ values
in a thread-safe manner.  The `ref` keyword allows you to _inline_ create a reference which you do not intend to
do anything else with.  This is really handy for passing _values_ to _reference parameters_ without needing to
also create a local variable.  It's absolutely _convenience only_ and not necessary for the long-term goals.

The `cleanup` keyword provides a way to defer a function until _the entire program shuts down._

### The `atlas` system

The atlas system should behave similarly to how it currently exists, but it should not enforce any specific
serialization type in the atlas file.  Instead, it should expect the following:

0. A typeless file called `atlas` next to the `go.mod` file
1. Some fields will be reserved for things like 'silent' and 'verbose' modes
2. The `atlas` package will provide thread-safe access to the content
3. The atlas file must be a known format, like JSON, XML, YAML, TOML, etc... with the reserved fields in the global space
4. The atlas file must be "live-watched" for updates in real time
5. It should gracefully be unable to parse it, record the anomaly, and continue
6. You would call `atlas.Parse[TOut](...func(any) TOut)` to then gain access to a populated object of your values.
7. If you provide a function to Parse, it should just call that to generate TOut (as you are defining the parsing function)

This means it should determine the atlas encoding scheme _on the fly_ and without _any_ configuration.  This system
should be as streamlined and polished as humanly possible because it _makes or breaks_ the _**entire**_ damn project!

Things like `atlas.ObservanceWindow` or `atlas.Precision` are specific to the system calling it.  As such,
each of those should be able to parse their fields from the atlas file without stepping on each other - so long as they
don't re-use the same field names.

NOTE: This means namespacing of fields is a concern we should consider addressing =)

    // JSON atlas Example
    {
        "core": {
            "verbose": true,
            "printPreamble": true,
        },
        "tiny": {
            "precision": 256,
        }
    }

---

    // XML atlas Example
    <core>
        <verbose>true</verbose>
        <printPreamble>true</printPreamble>
    </core>
    <tiny>
        <precision>256</precision>
    </tiny>

---

    // YAML atlas Example
    core:
        verbose: true
        printPreamble: true
    tiny:
        precision: 256

---

    // TOML atlas Example
    [core]
    verbose = true
    printPreamble = true

    [tiny]
    precision = 256

---

    // SPK atlas Example
    :Core:
        -verbose-true-
        -printPreamble-true-
    :
    :tiny:
        -precision-256-
    :

### Component Accessor Methods

Before we get into vectorization of a structure, I'd like to briefly discuss my thoughts on an efficient way to
identify, encapsulate, and access a _single component_ using the same _member name._  For example, let's say
you'd like to build an interface for anything "positionable" - you'd likely start with this:

    type Positionable interface {
        GetX() int
        SetX(int)
        GetY() int
        SetY(int)
    }

That's _quite clunky!_  So, let's unify the getting and setting into a single implicit operation:

    // A "fluently" positionable interface
    type Positionable interface {
        X(...int) int
        Y(...int) int
    }

    // An example implementation of "X"
    func (a Structure) X(val ...int) int {
        if len(val) > 0 {
            a.x = val[0]
        }
        return a.x
    }

This is what I call the "setter and/or getter" pattern, and it backs the entirety of the vectorization system because
it's _standardizable!_  What I mean is that using only a member name - `X` or `Y` - you **compile** the appropriate
path in your own code.  While that may not sound monumental, I assure you that significance will make far more sense
when we reach vectorization.  For now, let's consider a different question: why should a setter always return the
value?

Well, if X is an _overflowable_ value bound to a closed interval like `[0,42]` and you attempt to set it outside
of that range, it should gracefully return you the value it _actually_ set the component to!  You, the caller,
should have no knowledge of the logic driving that component of a structure's value except how to _**interface**_
with it.

For most code, the above pattern is _**absolutely more than adequate!**_  In fact, it's _preferred_ simply for
the fact that it facilitates _fluent chains!_  Most structures should be able to gracefully handle general
component access without _**forcing**_ complexity on _Other_.

But, as we introduce _vectorization,_ a new issue arises: what if the structure would like to report an _**error**_ 
occurred during access?  For example, what if you want to move a `Positionable` window structure, but you don't have
permission to do so?  In that case, this is the error variant of the setter and/or getter pattern:

    // A "gracefully" positionable interface
    type Positionable interface {
        X(...int) (int, error)
        Y(...int) (int, error)
    }

So, to recap:

0. Structures can have settable and/or gettable components identified by a _shared_ name
1. These can be accessed using either fluent or graceful accessor methods, as defined above

### Letter Vectors

Now, what if you have two structures - a `Primitive` one:

    type Primitive structure {
        X int
        Y int
    }

And an `Advanced` one:

    type Advanced structure {
        x int
        y int
    }

    func (a Advanced) X(...int) (int, error) {
        ...
        return a.x
    }

    func (a Advanced) Y(...int) (int, error) {
        ...
        return a.y
    }

You currently can't access `X` using the same pattern, as one is a field and the other a method.  But why would
a method not be able to implicitly handle both?  All we're attempting is standard getters and setters of the same
named member - a vector _component_ - but to do so would require reflection at runtime to execute a pretty basic
operation.

To handle this, GoOSE reserves three packages - `ltr`, `get`, and `set` - for self-describing vector type contracts.
If a type derives from one of those packages, it's considered to describe a fuzzy contract of what structures it can
interface with.  This allows your methods to articulate their intention by their method signature:

    // "I will move the X and Y components of any vector by the provided amount" 
    func Move[T any](amount T, vec ltr.XY[T, T]) (out ltr.XY[T, T], err error) {
        // Get the components individually
        x, err := vec.X()
        y, err := vec.Y()
    
        // Or get them as a set
        x, y, err = vec.Get()
    
        // Set the components individually
        x, err = vec.X(x+amouny)
        if err != nil {
            return nil, err
        }
        y, err = vec.Y(y+amount)
        if err != nil {
            return nil, err
        }

        // Or set them as a set
        x, y, err = vec.Set(x+amount, y+amount)
    
        return vec, nil
    }

A "vector," however, is an _abstract concept!_  There is no "type definition" backing it because it simply
articulates that the `Move` method takes any structure which has X and Y members either satisfying the 
_**"setter and/or getter" pattern**_ described above or as _**explicit fields.**_  Then, it exposes convenience
methods for directly "setting and/or getting" the components individually or collectively.

To call the method, simply pass in _any_ object that satisfies the contract - just as with interfaces!  Or,
to create an anonymous inline vector, you can either cast it from an object that satisfies the contract:

    var obj MyObj
    vec := ltr.XY[int, int](obj)

Or, you can do so with a vector literal:

    var x, y int
    
    vec := ltr.XY {
        x,
        y
    }

The type name separates components using _PascalCase,_ and the order/number of the members must match the name.
This allows you to perform execution directly on _qualities_ of a structure, which are similar to working with
"tuples" in other languages without explicitly _defining_ the tuple (minimizing clutter)

    // A "color space" vector
    var y, cb, cr float64
    vec := ltr.YCbCr {
        y,
        cb,
        cr
    }

    // A "physics" vector
    var f, m, a float64
    vec := ltr.FMA {
        f,
        m,
        a
    }

Above are type-inferred vectors, where the types are implicitly gleaned from the referenced local variables
and not included with the type name.  If including the type brackets, you must explicitly define the type of
_**all**_ vector components:

    var obj SomeOtherStruct
    vec := ltr.DeltaEpsilonGamma[TDelta int, TEpsilon func(), TGamma any] {
        42,
        nil,
        obj,
    }

The final aspect to consider as what I refer to as "wildcard" vectors:

    // "Wildcard" vector parameter and output types
    func Move(vec ltr.ABC) (out ltr.ABC, err error) {
        ...
        // Accepts ANY type satisfying the ABC component pattern
        // When getting, components yield their value as `any` type
        // When setting, you must introspect before using the appropriate type - otherwise, it will error at runtime
    }

    // "Hard typed"  vector parameter and output types
    func Move(vec ltr.ABC[any, any, any]) (out ltr.ABC[any, any, any], err error) {
        ...
        // Accepts only ltr.ABC[any, any any] typed objects
    }

    // Wildcards don't apply to literals, as their types cannot be gleaned at compile time
    var vecA ltr.ABC // INVALID CODE
    vecB := ltr.ABC { } // INVALID CODE

Wildcards are very important!  They shift the switch block necessary to execute on _any_ type from the calling
code to the receiving code - allowing them to own the _**entire_** process of working with ANY componentized data.
The only place wildcards are not allowed are off of the `set` package, as the receiver would have no way of knowing
what the appropriate type is to assign.

Lastly, the _three_ packages are used to encapsulate the contract's access:

0. `ltr` vectors are both settable _**and**_ gettable
1. `get` vectors are only gettable and will not compile code that sets the components
2. `set` vectors are only settable and will not compile code that gets the components

This allows the both sides of the contract a language level guarantee of how their members will be leveraged
_at compile time._  In addition, this lays the foundation for _empathy_ as defined by _The Treaty of Orlando_ =)