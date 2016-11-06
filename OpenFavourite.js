//
// Copyright (c) 2016 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2016-11-06
//

ObjC.import('stdlib')

var Transmit = Application('Transmit')
Transmit.includeStandardAdditions = true

// Return favourite for UID or null
function findFave(uid) {
    var faves = Transmit.favorites.whose({
        identifier: uid
    })
    if (faves.length == 0) {
        return null
    }
    return faves[0]
}

// Open favourite with UID
function openFave(uid) {
    var doc = null,
        fav = findFave(uid)

    if (fav === null) {
        console.log('Favourite not found: ' + uid)
        return false
    }
    
    doc = Transmit.Document().make()
    // Activate transmit if necessary
    if (!Transmit.frontmost()) {
        Transmit.activate()
    }
    // Open fave
    doc.currentTab.connect({to:fav})
    return true
}

function run(argv) {
    var uid = argv[0]
    if (!openFave(uid)) {
        $.exit(1)
    }
}