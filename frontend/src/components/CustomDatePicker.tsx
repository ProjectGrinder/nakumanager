import { useState, useRef, useEffect } from "react";

type DatePickerProps = {
  value: Date | null;
  onChange: (date: Date | null) => void;
};

const months = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

const getDaysInMonth = (year: number, month: number) =>
  new Date(year, month + 1, 0).getDate();

const formatDisplayDate = (date: Date) => {
  const day = String(date.getDate()).padStart(2, "0");
  const month = months[date.getMonth()];
  const year = String(date.getFullYear()).slice(-2);
  return `${day} ${month} ${year}`;
};

export default function CustomDatePicker(props: DatePickerProps) {
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [showCalendar, setShowCalendar] = useState(false);
  const [currentMonth, setCurrentMonth] = useState<number>(
    new Date().getMonth()
  );
  const [currentYear, setCurrentYear] = useState<number>(
    new Date().getFullYear()
  );
  const pickerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (pickerRef.current && !pickerRef.current.contains(e.target as Node)) {
        setShowCalendar(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const selectDate = (day: number) => {
    const date = new Date(currentYear, currentMonth, day);
    setSelectedDate(date);
    setShowCalendar(false);
  };

  const days = Array.from(
    { length: getDaysInMonth(currentYear, currentMonth) },
    (_, i) => i + 1
  );

  return (
    <div className="relative inline-block text-sm" ref={pickerRef}>
      <button
        onClick={() => setShowCalendar(!showCalendar)}
        className="px-3 py-1.5 rounded-lg bg-gray-700 text-gray-200 border-none"
      >
        {selectedDate ? formatDisplayDate(selectedDate) : "Select date"}
      </button>

      {showCalendar && (
        <div className="absolute left-0 z-10 mt-1 w-56 rounded border-none text-gray-400 bg-gray-700 shadow-lg">
          <div className="flex items-center justify-between px-4 py-2 border-b">
            <button
              className="hover:text-gray-300"
              onClick={() => {
                if (currentMonth === 0) {
                  setCurrentMonth(11);
                  setCurrentYear(currentYear - 1);
                } else {
                  setCurrentMonth(currentMonth - 1);
                }
              }}
            >
              ◀
            </button>
            <span className="font-medium">
              {months[currentMonth]} {currentYear}
            </span>
            <button
              className="hover:text-gray-300"
              onClick={() => {
                if (currentMonth === 11) {
                  setCurrentMonth(0);
                  setCurrentYear(currentYear + 1);
                } else {
                  setCurrentMonth(currentMonth + 1);
                }
              }}
            >
              ▶
            </button>
          </div>

          <div className="grid grid-cols-7 gap-1 p-2 text-center text-xs text-gray-500">
            {["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"].map((d) => (
              <div key={d}>{d}</div>
            ))}
          </div>

          <div className="grid grid-cols-7 gap-1 p-2 text-center">
            {Array(new Date(currentYear, currentMonth, 1).getDay())
              .fill(null)
              .map((_, i) => (
                <div key={`empty-${i}`}></div>
              ))}
            {days.map((day) => (
              <button
                key={day}
                onClick={() => selectDate(day)}
                className={`p-1 rounded hover:bg-gray-600 ${
                  selectedDate &&
                  selectedDate.getDate() === day &&
                  selectedDate.getMonth() === currentMonth &&
                  selectedDate.getFullYear() === currentYear
                    ? "bg-gray-600 text-white"
                    : ""
                }`}
              >
                {day}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
