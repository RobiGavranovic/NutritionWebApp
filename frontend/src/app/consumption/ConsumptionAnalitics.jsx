"use client";

import { PieChart } from "@mui/x-charts/PieChart";
import { BarChart } from "@mui/x-charts/BarChart";
import { useEffect, useState } from "react";

export default function ConsumptionAnalitics() {
  const [data, setData] = useState({ success: 0, failure: 0, noInput: 1 });
  const [range, setRange] = useState(7);

  const fetchData = (selectedRange) => {
    fetch(
      `http://localhost:3200/consumption/getConsumptionStatistics?range=${selectedRange}`,
      {
        method: "GET",
        credentials: "include",
      }
    )
      .then((res) => res.json())
      .then((res) => {
        if (
          typeof res.success === "number" &&
          typeof res.fail === "number" &&
          typeof res.noInput === "number" &&
          Array.isArray(res.intakeHistory)
        ) {
          setData({
            success: res.success,
            failure: res.fail,
            noInput: res.noInput,
            intakeHistory: res.intakeHistory,
          });
        }
      })
      .catch((err) => {
      });
  };

  useEffect(() => {
    fetchData(range);
  }, [range]);

  const successFailureTotal = data.success + data.failure;
  const total = successFailureTotal + data.noInput;

  const chartData =
    successFailureTotal > 0
      ? [
          { id: 0, value: data.success, label: "Success", color: "#16a34a" },
          { id: 1, value: data.failure, label: "Failure", color: "#dc2626" },
        ]
      : [{ id: 2, value: data.noInput, label: "No Input", color: "#6b7280" }];

  return (
    <div className="w-full flex flex-col items-center gap-6">
      {/* Centered Title */}
      <h2 className="text-2xl font-semibold text-gray-800 text-center">
        Consumption Analytics
      </h2>

      {/* Centered Range Selector */}
      <div className="flex items-center gap-2">
        <label className="text-sm font-medium text-gray-700">
          Select Range:
        </label>
        <select
          value={range}
          onChange={(e) => setRange(parseInt(e.target.value))}
          className="border border-gray-300 rounded p-2 text-sm"
        >
          {[7, 30, 60, 120, 210, 365].map((days) => (
            <option key={days} value={days}>
              {days} days
            </option>
          ))}
        </select>
      </div>

      {/* Two-column layout */}
      <div className="flex flex-col lg:flex-row justify-center items-stretch gap-10 w-full mt-4">
        {/* Left: Pie Chart + Legend */}
        <div className="w-full flex flex-col items-center gap-4">
          <PieChart
            series={[
              {
                data: chartData,
                innerRadius: 40,
                outerRadius: 100,
                paddingAngle: 3,
                cornerRadius: 3,
                cx: 150,
                cy: 100,
              },
            ]}
            width={300}
            height={220}
            colors={chartData.map((item) => item.color)}
            hideLegend={{ hidden: true }}
          />

          <div className="space-y-1">
            <div className="flex items-center space-x-2 text-sm text-gray-700">
              <span
                className="w-3 h-3 rounded-full inline-block"
                style={{ backgroundColor: "#16a34a" }}
              ></span>
              <span>
                Success:{" "}
                {successFailureTotal > 0
                  ? Math.round((data.success / successFailureTotal) * 100)
                  : 0}
                %
              </span>
            </div>
            <div className="flex items-center space-x-2 text-sm text-gray-700">
              <span
                className="w-3 h-3 rounded-full inline-block"
                style={{ backgroundColor: "#dc2626" }}
              ></span>
              <span>
                Failure:{" "}
                {successFailureTotal > 0
                  ? Math.round((data.failure / successFailureTotal) * 100)
                  : 0}
                %
              </span>
            </div>
            <div className="flex items-center space-x-2 text-sm text-gray-700">
              <span
                className="w-3 h-3 rounded-full inline-block"
                style={{ backgroundColor: "#6b7280" }}
              ></span>
              <span>
                No Input:{" "}
                {total > 0 ? Math.round((data.noInput / total) * 100) : 0}%
              </span>
            </div>
          </div>
        </div>

        {/* Consumption History Bar Chart */}
        <div className="w-full flex justify-center">
          <div className="w-full max-w-[500px]">
            <BarChart
              xAxis={[
                {
                  scaleType: "band",
                  data: (data.intakeHistory || []).map((entry) => entry.label),
                },
              ]}
              series={[
                {
                  data: (data.intakeHistory || []).map((entry) =>
                    Math.round(entry.avgCalories)
                  ),
                  color: "#3b82f6",
                },
              ]}
              height={220}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
