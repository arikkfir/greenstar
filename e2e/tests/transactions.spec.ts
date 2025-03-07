import { expect, Locator, test } from "@playwright/test"

test.describe("transactions", {}, () => {
    test("title", async ({ page }) => {
        await page.goto("/transactions")
        await expect(page).toHaveTitle(/Transactions - GreenSTAR/)
    })

    test("Root accounts shown", async ({ page }) => {
        await page.goto("/transactions")
        const nodesLocator: Locator =
                  page.locator("main ul[role=tree] li[role=treeitem] div.account-label p:nth-child(1)")
        await expect(nodesLocator).toHaveCount(7)
    })

    test("Default transactions shown", async ({ page }) => {
        await page.goto("/transactions")

        // Locators
        const txContainer: Locator = page.locator("main .transactions-container .transactions-grid-content")
        const txScrollerArea: Locator  = txContainer.locator(".MuiDataGrid-virtualScroller")
        const txRowsArea: Locator  = txContainer.locator(".MuiDataGrid-virtualScrollerContent")
        await expect(txScrollerArea).toHaveCount(1)

        // Keep scrolling down until we find all expected rows
        const referenceIDs: Set<string> = new Set<string>()
        while (true) {
            console.log(`**[ ${referenceIDs.size} rows ]*****************************************************`)
            const txRows = await txRowsArea.locator("div[role=row]").all()
            const txIDs  = txRows.map(async cell => await cell.getAttribute("data-id"))
            for (let i = 0; i < txIDs.length; i++) {
                const id = await txIDs[i]
                referenceIDs.add(id)
            }

            await txScrollerArea.evaluate(e => {
                console.log(`Scrolling down... (current scrollTop is ${e.scrollTop}, scrollHeight is ${e.scrollHeight})`)
                e.scrollBy(0, 300)
            })

            await page.waitForTimeout(1000)
            await expect(page.locator("span.MuiSkeleton-root")).toHaveCount(0)

            if (referenceIDs.size >= 199) {
                break
            }
        }

        expect(referenceIDs.size).toEqual(199)
    })
})
