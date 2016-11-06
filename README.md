Alfred Transmit
===============

Rapidly search Transmit favourites in [Alfred][alfredapp].


Usage
-----

- `ftp [<query>]` — View/filter favourites.
    - `↩` — Open favourite in Transmit.


Configuration
-------------

The workflow is pre-configured to use the `Favorites.xml` saved by the non-Mac App Store version of Transmit at `~/Library/Application Support/Transmit/Favorites/Favorites.xml`.

If you bought Transmit from the Mac App Store, you'll need to change the `FAVES_PATH` variable in the workflow configuration sheet to point to your `Favorites.xml` file.


Licencing & thanks
------------------

- [Transmit][transmit] goes without saying.
- The update icon is from the [Material Design Iconic Font][material] ([SIL licence and others][material-licence]) by [Sergey Kupletsky][sergey].
- The [AwGo library][awgo] ([MIT Licence][mit]) takes care of the workflowy stuff.


[alfredapp]: https://www.alfredapp.com/
[awgo]: https://godoc.org/gogs.deanishe.net/deanishe/awgo
[mit]: https://raw.githubusercontent.com/deanishe/alfred-ssh/master/LICENCE.txt
[octicons]: https://octicons.github.com/
[sil]: http://scripts.sil.org/OFL
[material]: https://zavoloklom.github.io/material-design-iconic-font/
[material-licence]: https://zavoloklom.github.io/material-design-iconic-font/license.html
[sergey]: http://twitter.com/zavoloklom
[transmit]: https://panic.com/transmit/
