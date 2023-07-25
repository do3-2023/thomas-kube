import { GetServerSideProps, InferGetServerSidePropsType } from "next";
import getConfig from "next/config";
import Link from "next/link";

const { serverRuntimeConfig } = getConfig();

export interface City {
  id: number;
  department_code: string;
  insee_code: string;
  zip_code: string;
  name: string;
  lat: number;
  lon: number;
}

export default function Cities({ city }: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <main className="flex flex-col items-center justify-between p-8 md:p-24">
      <h1 className="text-4xl md:text-6xl font-bold text-center mb-8">Welcome to my city website!</h1>
      <div className="text-center">
        <h2 className="text-xl md:text-2xl font-bold mb-4">Look at this super amazing city:</h2>
        <p className="text-lg md:text-xl mb-2">{city.name}</p>
        <p className="text-lg md:text-xl mb-2">{city.zip_code}</p>
        <p className="text-lg md:text-xl mb-2">{city.department_code}</p>
        <p className="text-lg md:text-xl mb-2">{city.insee_code}</p>
        <p className="text-lg md:text-xl mb-4">
          Latitude: {city.lat} | Longitude: {city.lon}
        </p>
      </div>

      <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded">
        <Link href="/cities">Click here to see another city</Link>
      </button>
    </main>
  );
}

export const getServerSideProps: GetServerSideProps<{ city: City }> = async () => {
  let city: any;
  try {
    const response = await fetch("http://backend-service.backend.svc.cluster.local/city/random");
    city = await response.json();
    console.log(city);
  } catch (error) {
    throw new Error(`Error fetching data to ${serverRuntimeConfig.baseAPI}: ${error}`);
  }
  return {
    props: {
      city,
    },
  };
};
