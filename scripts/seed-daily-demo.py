#!/usr/bin/env python3
"""Fill empty daily fields with labeled examples; existing classroom entries stay intact.
Run while the app is closed: python3 scripts/seed-daily-demo.py [data-directory]
"""
import copy
import datetime as dt
import json
import os
from pathlib import Path
import sys

root = Path(sys.argv[1]) if len(sys.argv) > 1 else Path(__file__).resolve().parents[1] / 'data'
config = json.loads((root / 'config.json').read_text())
path = root / 'data.json'
data = json.loads(path.read_text()) if path.exists() else {'version': 1, 'missing_homework': [], 'records': {}}
# Preserve an exact local copy so the seeded preview is reversible.
backup = root / 'data.before-demo.json'
if path.exists() and not backup.exists():
    backup.write_bytes(path.read_bytes())
    backup.chmod(0o600)
today = dt.date.today()
examples = [
    (['國語習作第 32～33 頁', '數學課本第 56 頁，完成第 1～5 題', '明天帶直尺、三角板及水壺'], '數學第四單元：分數的加減'),
    (['國語生字第七課，每字兩遍', '數學習作第 28～29 頁', '閱讀課外讀物 20 分鐘，記錄一句喜歡的句子', '家長簽閱聯絡簿'], '國語第七課：生字與語詞'),
    (['完成自然觀察紀錄', '英語單字 Unit 3，每字三遍'], '英語 Unit 3 單字小考'),
    (['整理本週學習單', '週末閱讀一本書，準備分享心得', '下週一帶美勞材料：色紙、膠水、剪刀'], ''),
    (['國語習作第 30～31 頁', '社會課本第 40～43 頁預習', '數學訂正並請家長簽名', '自然學習單完成第二部分', '明天穿運動服，攜帶跳繩'], '自然第二單元：水的變化'),
    (['完成數學練習卷', '準備兩分鐘閱讀分享'], '社會第二單元複習'),
    (['閱讀 20 分鐘', '檢查鉛筆盒與書包用品'], '國語第六課複習'),
]
for offset in range(7):
    day = today - dt.timedelta(days=offset)
    key = day.isoformat()
    record = data['records'].setdefault(key, {
        'date': key, 'updated_at': dt.datetime.now().astimezone().isoformat(),
        'duty': {'date': key, 'displayDate': f'{day.year}年{day.month}月{day.day}日', 'isWorkday': day.weekday() < 5 and key not in config.get('holidays', []), 'dutyStudents': [], 'lunchAssignments': []},
        'students': copy.deepcopy(config.get('students', [])), 'timetable': copy.deepcopy(config.get('timetable', [])),
        'period_times': config.get('settings', {}).get('period_times', []), 'missing_homework': [],
        'manual_duty': False, 'manual_lunch': False, 'layout_version': 1,
    })
    items, review = examples[offset]
    if not record.get('contact_items'):
        record['contact_items'] = ['【範例】' + item for item in items]
    if not record.get('next_day_review') and review:
        record['next_day_review'] = '【範例】' + review
for offset in range(7):
    day = today - dt.timedelta(days=offset)
    record = data['records'][day.isoformat()]
    previous = data['records'].get((day - dt.timedelta(days=1)).isoformat(), {})
    if not record.get('review_from_yesterday'):
        record['review_from_yesterday'] = previous.get('next_day_review', '')
temp = root / '.daily-demo.tmp'
temp.write_text(json.dumps(data, ensure_ascii=False, indent=2) + '\n')
temp.chmod(0o600)
os.replace(temp, path)
print(f'已補齊 {today - dt.timedelta(days=6)} ～ {today} 共七天範例，既有內容保留。')
print(f'原始資料保存在 {backup}；還原時請先關閉程式。')
