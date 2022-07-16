function delete_but(i) { //i == this

    let book = i.parentNode.parentNode.parentNode;
    // book.parentNode.parentNode.removeChild(book.parentNode) //js 删除通过父节点删除子节点
    let id = $("#" + book.children[1].id).attr("ind")    // 书的ID
    $.ajax({
        type: "post",
        url: "/admin/books/delete",
        data: {
            "id": id
        },
        success: (data) => {
            if (data.result == 0) {
                book.parentNode.parentNode.removeChild(book.parentNode)
            }
        }
    })
}

function BorrowBook(i) {
    // let book = i.parentNode.parentNode.parentNode;
    // book.parentNode.parentNode.removeChild(book.parentNode) //js 删除通过父节点删除子节点
    let id = i.parentNode.parentNode.attributes.ind.value;   // 书的ID

    // var id = i.parents("tr").find(".td.id").text();
    console.log(id);
    $.ajax({
        type: "post",
        url: "/user/Borrow",
        data: {
            "id": 3
        },
        success: (data) => {
            // if (data.result == 0) {
            //     book.parentNode.parentNode.removeChild(book.parentNode)
            // }
        }
    })
}

function ReturnBook(i) {
    // let book = i.parentNode.parentNode.parentNode;
    // book.parentNode.parentNode.removeChild(book.parentNode) //js 删除通过父节点删除子节点
    let id = i.parentNode.parentNode.attributes.ind.value;  // 书的ID
    console.log(id);
    $.ajax({
        type: "post",
        url: "/user/Return",
        data: {
            "id": id
        },
        success: (data) => {
            // if (data.result == 0) {
            //     book.parentNode.parentNode.removeChild(book.parentNode)
            // }
        }
    })
}


function bookDetails_but(i) { //i == this
    let book = i.parentNode.parentNode.parentNode;
    let id = $("#" + book.children[1].id).attr("ind")    // 书的ID
    console.log(id);
    // $.ajax({
    //     type: "get",
    //     url: "/admin/books/bookDetails",
    //     data: {
    //         "id": id
    //     },
    //     success: (data) => {
    //         console.log(data)
    //     }
    // })
    window.location.href = "/admin/books/bookDetails?id=" + id;
}