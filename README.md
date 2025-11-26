# Go Masker v2

[![build workflow](https://github.com/ggwhite/go-masker/actions/workflows/go.yml/badge.svg)](https://github.com/ggwhite/go-masker/actions)
[![GoDoc](https://godoc.org/github.com/ggwhite/go-masker?status.svg)](https://godoc.org/github.com/ggwhite/go-masker)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggwhite/go-masker)](https://goreportcard.com/report/github.com/ggwhite/go-masker)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/ggwhite/go-masker/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/ggwhite/go-masker.svg?style=flat-square)](https://github.com/ggwhite/go-masker/releases/latest)

Go Masker v2 is a tool for masking sensitive data in Go code. It provides a simple and convenient way to replace sensitive information, such as passwords or API keys, with placeholder values.

* [Install](#install)
* [Usage](#usage)
* [Abuse Masker](#abuse-masker)
* [Custom Masker](#custom-masker)
* [Contributors](#contributors)

## Install

To install Go Masker v2, you can use the following command:

```bash
go get -u github.com/ggwhite/go-masker/v2
```

## Usage

To use Go Masker v2, you can create a new instance of the `Masker` type and then use its methods to mask sensitive data. For example:

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
)

type Foo struct {
    Name   string `mask:"name"`
    Mobile string `mask:"mobile"`
}

func main() {
    foo := &Foo{
        Name:   "ggwhite",
        Mobile: "0987987987",
    }

    m := masker.NewMaskerMarshaler()

    t, err := m.Struct(foo)
    log.Println(t)
    log.Println(err)
}
```

This will produce the following output:

```
t = &{g**hite 0987***987}
err = <nil>
```

For more information about how to use Go Masker v2, please refer to the [documentation](https://pkg.go.dev/github.com/ggwhite/go-masker/v2).

## Abuse Masker

The abuse masker uses a trie data structure for efficient word matching and masking. It allows you to mask abusive or inappropriate words in text content, which is crucial for maintaining appropriate communication in chat applications, forums, and social media platforms.

### Basic Usage

```go
package main

import (
    "fmt"
    masker "github.com/ggwhite/go-masker/v2"
)

func main() {
    // Create abuse masker with predefined words
    abuseWords := []string{"bad", "terrible", "awful"}
    abuseMasker := masker.NewAbuseMaskerWithWords(abuseWords)
    
    text := "This is a bad word and terrible situation"
    maskedText := abuseMasker.Marshal("*", text)
    
    fmt.Printf("Original: %s\n", text)
    fmt.Printf("Masked:   %s\n", maskedText)
    // Output: This is a *** word and ******** situation
}
```

### Loading Words from File

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
)

func main() {
    loader := masker.NewAbuseWordLoader()
    
    // Load words from a file
    words, err := loader.LoadFromFile("abuse_words.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    abuseMasker := masker.NewAbuseMaskerWithWords(words)
    // Use the masker...
}
```

### Loading Words from S3

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)
//Your own function to load words from S3 and pass words in abuse masker
func loadWordsFromS3(bucket, key string) ([]string, error) {
    sess := session.Must(session.NewSession())
    svc := s3.New(sess)
    
    result, err := svc.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, err
    }
    defer result.Body.Close()
    
    loader := masker.NewAbuseWordLoader()
    return loader.LoadFromReader(result.Body)
}

func main() {
    words, err := loadWordsFromS3("my-bucket", "abuse-words.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    abuseMasker := masker.NewAbuseMaskerWithWords(words)
    // Use the masker...
}
```

### Using with Struct Tags

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
)

type User struct {
    Name    string `mask:"name"`
    Comment string `mask:"abuse"`
}

func main() {
    user := &User{
        Name:    "John Doe",
        Comment: "This is a bad comment",
    }
    
    m := masker.NewMaskerMarshaler()
    
    // Register abuse masker with words
    abuseMasker := masker.NewAbuseMaskerWithWords([]string{"bad", "terrible"})
    m.Register(masker.MaskerTypeAbuse, abuseMasker)
    
    maskedUser, err := m.Struct(user)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Masked: %+v", maskedUser)
}
```

## Custom Masker

You can also create a custom masker by implementing the `Masker` interface. For example:

```go
package main

import (
    "log"
    masker "github.com/ggwhite/go-masker/v2"
)

type MyEmailMasker struct{}

func (m *MyEmailMasker) Marshal(s, i string) string {
	return "myemailmasker"
}

type MyMasker struct{}

func (m *MyMasker) Marshal(s, i string) string {
	return "mymasker"
}

func main() {
    m := masker.NewMaskerMarshaler()

    // Register custom masker and override default masker
	m.Register(masker.MaskerTypeEmail, &MyEmailMasker{})

	log.Println(m.Marshal(masker.MaskerTypeEmail, "email")) // myemailmasker <nil>

	// Register custom masker and use it
	m.Register("mymasker", &MyMasker{})

	log.Println(m.Marshal("mymasker", "1234567")) // mymasker <nil>
}
```


## Contributors

Thanks to all the people who already contributed!

<a href="https://github.com/ggwhite/go-masker/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=ggwhite/go-masker" />
</a>
