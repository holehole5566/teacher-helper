// Run with Playwright installed: PLAYWRIGHT_MODULE=/path/to/playwright/index.mjs node scripts/capture.mjs
// Start frontend Vite first. Only synthetic data is injected; no Go backend or config.json is used.
import { mkdir } from 'node:fs/promises';
const { chromium } = await import(process.env.PLAYWRIGHT_MODULE || 'playwright');
const out = new URL('../public/screenshots/', import.meta.url).pathname;
await mkdir(out, { recursive: true });
const browser = await chromium.launch({ executablePath: process.env.CHROME_PATH || '/usr/bin/google-chrome', headless: true, args: ['--no-sandbox'] });
const page = await browser.newPage({ viewport: { width: 1280, height: 860 }, deviceScaleFactor: 1, timezoneId: 'Asia/Taipei' });
const errors = [];
page.on('pageerror', e => errors.push(e.message));
await page.clock.install({ time: new Date('2026-09-21T08:40:00+08:00') });
await page.clock.pauseAt(new Date('2026-09-21T08:40:00+08:00'));
await page.addInitScript(() => {
  const students = ['陳小宇','林語晴','王柏安','張雅涵','李品樂','黃子謙','吳思妤','劉宥辰','蔡欣宜','楊承恩','許詠晴','鄭睿哲'].map((name, i) => ({ name, seat_number: i + 1, duty_enabled: i !== 10, lunch_enabled: i !== 11 }));
  const settings = { semester_start_date: '2026-08-31', duty_group_size: 2, lunch_group_size: 5, duty_start_number: 1, lunch_start_number: 3, meal_buckets: ['飯桶','主菜','配菜','青菜','湯桶'], period_times: ['08:40','09:30','10:30','11:20','12:30','13:30','14:20','15:10'], countdown_times: ['08:40','09:30','10:30','11:20','13:30','14:20','15:10'], countdown_volume: 0.5, countdown_musics: [], countdown_time_music_map: [], audio_output_device: 'default', discord_webhook: '', auto_start: false };
  const timetable = [
    ['國語','數學','英語','自然','午休','體育','社會','綜合'],
    ['數學','國語','自然','美勞','午休','美勞','英語','彈性'],
    ['國語','數學','音樂','社會','午休','','',''],
    ['英語','國語','數學','自然','午休','體育','綜合','閱讀'],
    ['國語','數學','社會','健康','午休','音樂','綜合','班會'],
  ];
  const api = {
    HasPassword: () => true,
    GetStudents: () => students,
    GetSettings: () => settings,
    GetTimetable: () => timetable,
    GetMissingHomework: () => [{ subject: '數學', students: [7,12], note: '習作第 12–13 頁，明早補交' }, { subject: '國語', students: [3,9], note: '第三課生字練習' }],
    GetHolidays: () => ['2026-09-25','2026-09-28','2026-10-09','2026-10-26'],
    GetTodayDuty: () => ({ date: '2026-09-21', displayDate: '2026 年 9 月 21 日（星期一）', isWorkday: true, dutyStudents: students.slice(0,2), lunchAssignments: settings.meal_buckets.map((bucket,i) => ({ bucket, student: students[i+2] })) }),
    SetFullscreen: () => {}, DebugLog: () => {},
    ReportError: message => { throw new Error(message); },
    GetActiveCountdownMusicData: () => '',
  };
  window.go = { main: { App: Object.fromEntries(Object.entries(api).map(([name, fn]) => [name, async (...args) => structuredClone(fn(...args))])) } };
  window.demoEvents = {};
  window.runtime = { EventsOnMultiple: (name, callback) => { window.demoEvents[name] = callback; return () => {}; }, EventsOff: () => {} };
});
try {
  await page.goto(process.env.APP_URL || 'http://127.0.0.1:5173');
  await page.getByText('陳小宇', { exact: false }).first().waitFor();
  async function shot(name) {
    await page.evaluate(() => document.fonts.ready);
    await page.screenshot({ path: `${out}${name}.png`, animations: 'disabled' });
    console.log(`Captured ${name}`);
  }
  await shot('overview');
  for (const [label, name, ready] of [['課表設定','timetable','儲存課表'],['學生管理','students','新增'],['作業未交','homework','作業未交管理'],['設定','settings','儲存設定'],['假期管理','holidays','同步政府假日']]) {
    await page.getByRole('button', { name: label, exact: true }).click();
    await page.getByText(ready, { exact: true }).first().waitFor();
    await shot(name);
  }
  await page.setViewportSize({ width: 1600, height: 1200 });
  await page.getByRole('button', { name: '展示模式', exact: true }).click();
  await page.locator('.display .class-row').first().waitFor();
  await shot('display');
  await page.evaluate(() => window.demoEvents['countdown-trigger']('08:40'));
  await page.locator('.countdown-container').waitFor();
  await shot('countdown');
  if (errors.length) throw new Error(errors.join('\n'));
} finally { await browser.close(); }
