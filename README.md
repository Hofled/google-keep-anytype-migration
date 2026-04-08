# Google Keep To Anytype Migration
This tool is used for migrating [Google Keep](https://keep.google.com/) notes to [Anytype](https://anytype.io/).

The migration relies on the data format from using the [Google Takeout](https://takeout.google.com/) export tool.

## Roadmap
- [ ] Migration using [Anytype API](https://developers.anytype.io/docs/reference/2025-11-08/anytype-api/), Google Keep Note ➡️ Anytype Page
    - [ ] Only text content
    - [ ] Checkbox list formatting
    - [ ] Note annotations ➡️ Anytype Bookmarks
    - [ ] Google Keep Labels ➡️ Anytype Tags
    - [ ] All text formatting
    - [ ] Include images once the [file uploading API](https://github.com/anyproto/anytype-heart/pull/2843) is supported
