# blogger

This is an example Go application to demonstrate how first-class functions in
Go can be utilised for idiomatic query building, and relationship loading.

To read more about this see the following,

* [ORMs and Query Building in Go](https://andrewpillar.com/programming/2019/07/13/orms-and-query-building-in-go)
* [Working with SQL Relations in Go - Part 1](https://andrewpillar.com/programming/2020/04/07/working-with-sql-relations-in-go-part-1/)
* [Working with SQL Relations in Go - Part 2](https://andrewpillar.com/programming/2020/04/07/working-with-sql-relations-in-go-part-2/)
* [Working with SQL Relations in Go - Part 3](https://andrewpillar.com/programming/2020/04/07/working-with-sql-relations-in-go-part-3/)
* [Working with SQL Relations in Go - Part 4](https://andrewpillar.com/programming/2020/04/07/working-with-sql-relations-in-go-part-4/)
* [Working with SQL Relations in Go - Part 5](https://andrewpillar.com/programming/2020/04/07/working-with-sql-relations-in-go-part-5/)

To build and run this application simply clone, and run the `make.sh` script.

    $ git clone https://github.com/andrewpillar/blogger
    $ cd blogger
    blogger $ ./make.sh
    blogger $ ./blogger.out

Then via cURL you can start querying the API.

    $ curl -s "localhost:8080/categories/5"
    $ curl -s "localhost:8080/posts?search=rev"
    $ curl -s "localhost:8080/posts?tag=hosts"
