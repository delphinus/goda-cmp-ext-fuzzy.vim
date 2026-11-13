-- Types the demo into demo/sample/main.go, then quits (GODA_DEMO_DRIVE=1).

local steps = {
  { 1500, "G" },
  { 300, "O" },
}

local function input(text)
  for c in text:gmatch "." do
    table.insert(steps, { 180, c })
  end
end

-- Println out of "prl", TrimSpace out of "tsp" and sentence out of "snt".
input "fmt.prl"
vim.list_extend(steps, { { 1800, "<C-n>" }, { 1500, "<C-y>" }, { 300, "(" } })
input "strings.tsp"
vim.list_extend(steps, { { 1800, "<C-n>" }, { 1500, "<C-y>" }, { 300, "(" } })
input "snt"
vim.list_extend(steps, {
  { 1800, "<C-n>" },
  { 1500, "<C-y>" },
  { 300, ")" },
  { 150, ")" },
  { 1000, "<Esc>" },
})

local function feed(keys)
  vim.api.nvim_feedkeys(vim.api.nvim_replace_termcodes(keys, true, false, true), "t", false)
end

local function run(i)
  if i > #steps then
    vim.defer_fn(function()
      if vim.env.GODA_DEMO_LOG then
        vim.fn.writefile(vim.fn["goda#log"](), vim.env.GODA_DEMO_LOG)
      end
      vim.cmd "qa!"
    end, 1500)
    return
  end
  vim.defer_fn(function()
    feed(steps[i][2])
    run(i + 1)
  end, steps[i][1])
end

-- Start once gopls has attached.
vim.wait(10000, function()
  return #vim.lsp.get_clients { bufnr = 0, name = "gopls" } > 0
end, 100)
-- Clear what startup said (the recording terminal answers no queries).
vim.cmd "echo '' | redraw"
run(1)
