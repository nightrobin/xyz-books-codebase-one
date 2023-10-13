# XYZ Books Codebase 1

---

 This program will cater the CRUD UI and API (backend) for the XYZ Books. 

---


## Instructions
1. Create a database in MySQL.
2. Edit and save the **.env** file that is included in the codebase 1.
	- Input the correct database credentials in the DB fields.
	- Input the correct API Host and API Port (default is recommended).
3. Go to the project directory via terminal.
4. Run `go get`
5. Run `go build` 
6. After a successful build, you may run the `main` executable.

## Scope & Limitations
- As Frontend Frameworks are **not allowed** to be used, HTML Templates were used instead to handle the Client UI for the codebase 1.
- For the requirement ***"Every specific record page/endpoint should be given its ISBN-13 in the URL path."*** It create an issue if the ISBN-13 is not yet set in a particular book. To properly handle this, viewing of that record is only allowed in the CRUD UI: Books Index and in the API: Books - Get All.
- The specific fetching and viewing of the book will be allowed, once that book has been set a valid ISBN-13.
- Given the huge datasets, I implemented a pagination for the CRUD UI and API.
- For database, I use MySQL for this project.

## Versions used
- Go: 1.21.1
- MySQL: 8

## CRUD UI Reference
**URL:** http://localhost:3000

### CRUD UI Search
To search in the CRUD UI, use the provided search input.

Take note that it will only consider the following criterias when searching:
- **Books Page:** Title, ISBN-13, ISBN-10, Author, Publication Year, Publisher Name
- **Authors Page:** First Name, Middle Name, Last Name
- **Publisher Page:** Name

## API Endpoints Reference
|AUTHOR|Get All  |
|--|--|
|URL  |http://localhost:3000/api/authors  |
|Method  |GET  |

**Query string parameters**
Example URL: http://localhost:3000/api/authors?keyword=john&limit=10000&page=2
| Field | Specification | Description |
|----------|----------|----------|
| keyword | optional | Keyword for searching. It uses the same criterias like the CRUD UI Author search.  | 
| limit | optional | Record limit to fetch.  |
| page | optional | The offset page, only works if the **limit** query string parameter is also defined.  |

**Sample Response**

    {
    "message": "Successfully retrieved the authors.",
    "count": 9,
    "page": 1,
    "data": [
        {
            "ID": 1,
            "FirstName": "Joel",
            "MiddleName": "",
            "LastName": "Hartse"
        },
        {
            "ID": 2,
            "FirstName": "Hannah",
            "MiddleName": "P.",
            "LastName": "Templer"
        },
        {
            "ID": 3,
            "FirstName": "Marguerite",
            "MiddleName": "Z.",
            "LastName": "Duras"
        },
        {
            "ID": 4,
            "FirstName": "Kingsley",
            "MiddleName": "",
            "LastName": "Amis"
        },
        {
            "ID": 5,
            "FirstName": "Fannie",
            "MiddleName": "Peters",
            "LastName": "Flagg"
        },
        {
            "ID": 6,
            "FirstName": "Camille",
            "MiddleName": "Byron",
            "LastName": "Paglia"
        },
        {
            "ID": 7,
            "FirstName": "Rainer",
            "MiddleName": "Steel",
            "LastName": "Rilke"
        },
        {
            "ID": 8,
            "FirstName": "Robinson",
            "MiddleName": "Erlano",
            "LastName": "Joaquin"
        },
        {
            "ID": 9,
            "FirstName": "Robinson",
            "MiddleName": "Erlano",
            "LastName": "Joaquin"
        }
    ],
    "errors": null
    }
---
|AUTHOR|Get One  |
|--|--|
|URL  |http://localhost:3000/api/authors/:id  |
|Method  |GET  |

**Query Params:**
| Field | Specification | Description |
|----------|----------|----------|
| id | integer, required | Unique identifier of an author. |

**Sample Response**

    {
    "message": "Successfully retrieved the author.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "FirstName": "Joel",
        "MiddleName": "",
        "LastName": "Hartse"
    },
    "errors": null
    }
---
|AUTHOR|Create  |
|--|--|
|URL  |http://localhost:3000/api/authors  |
|Method  |POST  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| FirstName | string, required, max length: 255 | First name of the author. |
| MiddleName | string, optional, max length: 255 | Middle name of the author. |
| LastName | string, required, max length: 255 | Last name of the author. |

**Sample Request Body**

    {
    "FirstName":  "Robinson",
    "MiddleName":  "Erlano",
    "LastName":  "Joaquin"
    }

**Sample Response**

    {
    "message": "Successfully added the Author",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 10,
        "FirstName": "Robinson",
        "MiddleName": "Erlano",
        "LastName": "Joaquin"
    },
    "errors": null
    }
---
|AUTHOR|Update  |
|--|--|
|URL  |http://localhost:3000/api/authors/:id  |
|Method  |PATCH  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| id | integer, required | Unique identifier of an author. |
| FirstName | string, required, max length: 255 | First name of the author. |
| MiddleName | string, optional, max length: 255 | Middle name of the author. |
| LastName | string, required, max length: 255 | Last name of the author. |

**Sample Request Body**

    {
    "FirstName":  "Robinson",
    "MiddleName":  "Erlano",
    "LastName":  "Joaquin"
    }

**Sample Response**

    {
    "message": "Successfully updated the author.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "FirstName": "John",
        "MiddleName": "B",
        "LastName": "Smith"
    },
    "errors": null
    }
---
|AUTHOR|Delete  |
|--|--|
|URL  |http://localhost:3000/api/authors/:id  |
|Method  |DELETE  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| id | integer, required | Unique identifier of an author. |

**Sample Response**

    {
    "message": "Successfully deleted the author",
    "count": 1,
    "page": 1,
    "data": null,
    "errors": null
    }
---
|BOOK|Get All  |
|--|--|
|URL  |http://localhost:3000/api/books  |
|Method  |GET  |

**Query string parameters**
Example URL: http://localhost:3000/api/books?keyword=bear&limit=10000&page=2
| Field | Specification | Description |
|----------|----------|----------|
| keyword | optional | Keyword for searching. It uses the same criterias like the CRUD UI Book search.  | 
| limit | optional | Record limit to fetch.  |
| page | optional | The offset page, only works if the **limit** query string parameter is also defined.  |

**Sample Response**

    {
    "message": "Successfully retrieved the books.",
    "count": 4,
    "page": 1,
    "data": [
        {
            "ID": 7,
            "Title": "Book 7",
            "Isbn13": "",
            "Isbn10": "",
            "PublicationYear": 1231,
            "PublisherID": 1,
            "ImageURL": "",
            "Edition": "",
            "ListPrice": 123,
            "AuthorIDs": "[3]"
        },
        {
            "ID": 8,
            "Title": "Book 9",
            "Isbn13": "",
            "Isbn10": "1891830856",
            "PublicationYear": 1312,
            "PublisherID": 1,
            "ImageURL": "",
            "Edition": "",
            "ListPrice": 12321,
            "AuthorIDs": "[5]"
        },
        {
            "ID": 9,
            "Title": "123123",
            "Isbn13": "9781603093989",
            "Isbn10": "",
            "PublicationYear": 1231,
            "PublisherID": 1,
            "ImageURL": "",
            "Edition": "",
            "ListPrice": 1232,
            "AuthorIDs": "[3]"
        },
        {
            "ID": 11,
            "Title": "123123",
            "Isbn13": "9781891830020",
            "Isbn10": "",
            "PublicationYear": 1312,
            "PublisherID": 1,
            "ImageURL": "",
            "Edition": "",
            "ListPrice": 123,
            "AuthorIDs": "[3]"
        }
    ],
    "errors": null
    }

---
|BOOK|Get One  |
|--|--|
|URL  |http://localhost:3000/api/books/:isbn  |
|Method  |GET  |

**Query Params:**
| Field | Specification | Description |
|----------|----------|----------|
| isbn | string, required, max length = 13 | Valid ISBN 13 of the book. |

**Sample Response**

    {
    "message": "Successfully retrieved the book.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 4,
        "Title": "Hey, Mister (Vol 1)",
        "Isbn13": "9781891830020",
        "Isbn10": "1891830023",
        "PublicationYear": 2000,
        "PublisherID": 3,
        "ImageURL": "",
        "Edition": "After School Special",
        "ListPrice": 1200,
        "AuthorIDs": "[2,5,6]"
    },
    "errors": null
    }
---
|BOOK|Create  |
|--|--|
|URL  |http://localhost:3000/api/books  |
|Method  |POST  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| Title | string, required, max length: 255 | Title of the book. |
| Isbn13 | string, optional, min length: 13, max length: 13 | A valid ISBN-13 of the book. |
| Isbn10 | string, optional, min length: 10 | A valid ISBN-10 of the book. |
| PublicationYear | string, required | The publication year of the book. |
| PublisherID | integer, required | A valid ID of the Publisher of the book. |
| ImageURL | string, required, max length: 255 | Direct Image URL of the Book. |
| Edition | string, required, max length: 255 | Edition of the book. |
| ListPrice | float, required | Price of the book. |
| AuthorIDs | array of integers, required | IDs of the Author(s). Atleast one valid author ID is required. |
**Sample Request Body**

    {
    "Title": "Book 1",
    "Isbn13": "9781603095020",
    "Isbn10": "048665088X",
    "PublicationYear": 2001,
    "PublisherID": 1,
    "ImageURL": "https://image.com/book.jpg",
    "Edition": "Book 0",
    "ListPrice": 123.12,
    "AuthorIDs": [2,3]
    }
   

**Sample Response**

    {
    "message": "Successfully added the book.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "Title": "Book 1",
        "Isbn13": "9781603095020",
        "Isbn10": "048665088X",
        "PublicationYear": 2001,
        "PublisherID": 1,
        "ImageURL": "https://image.com/book.jpg",
        "Edition": "Book 0",
        "ListPrice": 123.12,
        "AuthorIDs": [
            2,
            3
        ]
    },
    "errors": null
    }
---
|BOOK|Update  |
|--|--|
|URL  |http://localhost:3000/api/books/:id  |
|Method  |PATCH  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| Title | string, required, max length: 255 | Title of the book. |
| Isbn13 | string, optional, min length: 13, max length: 13 | A valid ISBN-13 of the book. |
| Isbn10 | string, optional, min length: 10 | A valid ISBN-10 of the book. |
| PublicationYear | string, required | The publication year of the book. |
| PublisherID | integer, required | A valid ID of the Publisher of the book. |
| ImageURL | string, required, max length: 255 | Direct Image URL of the Book. |
| Edition | string, required, max length: 255 | Edition of the book. |
| ListPrice | float, required | Price of the book. |
| AuthorIDs | array of integers, required | IDs of the Author(s). Atleast one valid author ID is required. |

**Sample Request Body**

    {
    "Title": "Book 2",
    "Isbn13": "9781603095020",
    "Isbn10": "048665088X",
    "PublicationYear": 2000,
    "PublisherID": 1,
    "ImageURL": "https://image.com/book2.jpg",
    "Edition": "Book 2",
    "ListPrice": 123.12,
    "AuthorIDs": [2,3]
    }

**Sample Response**

    {
    "message": "Successfully updated the book.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "Title": "Book 2",
        "Isbn13": "9781891830853",
        "Isbn10": "1891830023",
        "PublicationYear": 2000,
        "PublisherID": 1,
        "ImageURL": "https://image.com/book2.jpg",
        "Edition": "Book 2",
        "ListPrice": 123.12,
        "AuthorIDs": [
            2,
            3
        ]
    },
    "errors": null
    }
---
|BOOK|Delete  |
|--|--|
|URL  |http://localhost:3000/api/books/:id  |
|Method  |DELETE  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| id | integer, required | Unique identifier of the book. |

**Sample Response**

    {
    "message": "Successfully deleted the book",
    "count": 1,
    "page": 1,
    "data": null,
    "errors": null
    }
---
|PUBLISHER|Get All  |
|--|--|
|URL  |http://localhost:3000/api/publishers  |
|Method  |GET  |

**Query string parameters**
Example URL: http://localhost:3000/api/publishers?keyword=bear&limit=10000&page=2
| Field | Specification | Description |
|----------|----------|----------|
| keyword | optional | Keyword for searching. It uses the same criterias like the CRUD UI Publisher search.  | 
| limit | optional | Record limit to fetch.  |
| page | optional | The offset page, only works if the **limit** query string parameter is also defined.  |

**Sample Response**

    {
    "message": "Successfully retrieved the publishers.",
    "count": 4,
    "page": 1,
    "data": [
        {
            "ID": 1,
            "Name": "Paste Magazine"
        },
        {
            "ID": 2,
            "Name": "Publishers Weekly"
        },
        {
            "ID": 3,
            "Name": "Graywolf Press"
        },
        {
            "ID": 4,
            "Name": "McSweeney's"
        }
    ],
    "errors": null
    }

---
|PUBLISHER|Get One  |
|--|--|
|URL  |http://localhost:3000/api/publishers/:id  |
|Method  |GET  |

**Query Params:**
| Field | Specification | Description |
|----------|----------|----------|
| id | integer, required | Unique identifier of the publisher. |

**Sample Response**

    {
    "message": "Successfully retrieved the publisher.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "Name": "Paste Magazine"
    },
    "errors": null
        
    }
---
|PUBLISHER|Create  |
|--|--|
|URL  |http://localhost:3000/api/publishers  |
|Method  |POST  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| Name | string, required, max length: 255 | Name of the publisher. |
**Sample Request Body**

    {
    "Name": "International Bookstore"
    }

**Sample Response**

    {
    "message": "Successfully added the publisher.",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "Name": "International Bookstore"
    },
    "errors": null
    }
---
|PUBLISHER|Update  |
|--|--|
|URL  |http://localhost:3000/api/publishers/:id  |
|Method  |PATCH  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| Name | string, required, max length: 255 | Name of the publisher. |

**Sample Request Body**

    {
    "Name": "Local Bookstore"
    }

**Sample Response**

    {
    "message": "Successfully updated the publisher",
    "count": 1,
    "page": 1,
    "data": {
        "ID": 1,
        "Name": "Local Bookstore"
    },
    "errors": null
    }
---
|PUBLISHER|Delete  |
|--|--|
|URL  |http://localhost:3000/api/publishers/:id  |
|Method  |DELETE  |

**Request Body:**
| Field | Specification | Description |
|----------|----------|----------|
| id | integer, required | Unique identifier of the publisher. |

**Sample Response**

    {
    "message": "Successfully deleted the publisher",
    "count": 1,
    "page": 1,
    "data": null,
    "errors": null
    }