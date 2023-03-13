let next = document.getElementById("nextPageBtn")
next.addEventListener("click", function() {
    let currentURL = window.location.href
    let url = new URL(currentURL)
    if (url.searchParams.has("offset")) {
        let offset = parseInt(url.searchParams.get("offset"))
        url.searchParams.set("offset", offset + 20)
    } else {
        url.searchParams.set("offset", 20)
    }
    window.location.href = url.toString()
})

let previous = document.getElementById("previousPageBtn")
previous.addEventListener("click", function() {
    let currentURL = window.location.href
    let url = new URL(currentURL)
    if (url.searchParams.has("offset")) {
        let offset = parseInt(url.searchParams.get("offset"))
        if (offset - 20 > 0) {
            url.searchParams.set("offset", offset - 20)
        } else {
            url.searchParams.delete("offset")
        }
        window.location.href = url.toString()
    }
})