# go-tlv

TLV Parser Library

- fields are non-unique
  ```yaml
  message:
    obj1:
      field1: a
      field1: b
    obj1:
      field2: a
      field2: b
  ```

- the parser supports multiple root level fields:
  ```yaml
  first-message:
    field: value
  second-message:
    field: value
  ```
