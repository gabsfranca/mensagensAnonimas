import { SolidApexCharts } from "solid-apexcharts";
import { Component, createResource } from 'solid-js';

const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

interface TagStat {
    tag: string;
    count: number;
}

const fetchTagStats  = async (): Promise <TagStat[]> => {
    const res = await fetch(`${URL}/messages/tags/stats`, {
        credentials: 'include',
    });
    if (!res.ok) throw new Error('erro ao carregar estatisticas');
    return res.json();
};

const TagsBarChart: Component = () => {
    const [data] = createResource(fetchTagStats);

    return (
    <div>
      <h2>Mensagens por Tag</h2>
      <SolidApexCharts
        type="bar"
        options={{
          chart: {
            id: "tags-bar",
          },
          xaxis: {
            categories: data()?.map((d) => String(d.tag)) || [],
            labels: {
              rotate: -45,
              style: {
                fontSize: '12px'
              },
            },
          }
        }}
        series={[
          {
            name: "Quantidade",
            data: data()?.map((d) => d.count) || [],
          },
        ]}
        width="100%"
        height="400"
      />
    </div>
  );
};

export default TagsBarChart;