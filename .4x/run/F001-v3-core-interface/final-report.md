# Final Report — F001-v3-core-interface: F001 v3 core: 新 Masker interface + package-level 短函式

## Status
ready-for-review

## Summary
在新建的 `v3/` module 內完成核心基礎建設：新 `Masker` interface（`Mask(value string) string`）、functional option `WithMaskChar`、精簡 `MaskerType` 常數（`TypeMobile` 等）、package-level 便利函式、`MustMarshal`，並將全部 12 種 masker（含 generic first-N/last-N）遷移至新介面。遮罩輸出與 v2 逐字一致，零外部依賴，16/16 acceptance criteria 全數通過，`go test -race` 全綠。

## Open Issues
None — all issues resolved.

備註（非阻擋，已記於 review-report 與 test-report）：coverage 63.9%，未覆蓋部分集中於 `Struct()` 的 slice/map/ptr reflect 容器分支——該重構為獨立 feature（reflect cache，out of scope），本 feature 不設覆蓋率門檻，核心 interface / 各 masker / convenience / Option 路徑皆已覆蓋，不構成 open issue。
