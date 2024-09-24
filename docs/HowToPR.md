# Pull Request Template

### Description
A brief summary of what this pull request does and why it's necessary.

### Changes
A list of changes made in this pull request.

### The detail of the pull request.
Write the details here (e.g. affected APIs, updated code flow, etc).

### Action
- [ ] Bug fixing
- [ ] Enhancement
- [ ] Feature

### Review/Test checklist
- [ ] All unit tests pass
- [ ] Manually tested API endpoints
- [ ] Documentation has been updated to reflect any changes

# Pull Request Example

### Description
add is_publish query on get news api

### Changes
- add is_publish is true on get news query
- add is_publish is true on get news by code query

### The detail of the pull request.
- api :
 - /v1/news/{code}
 - /v1/news
- no changes on parameter and response format

### Action
- [ ] Bug fixing
- [x] Enhancement
- [ ] Feature

### Review/Test checklist
- [x] All unit tests pass
- [x] Manually tested API endpoints
- [ ] Documentation has been updated to reflect any changes