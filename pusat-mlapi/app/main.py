from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from pandas import json_normalize 
import pandas as pd
import requests
import pickle
from sklearn.preprocessing import LabelEncoder
import numpy as np
import joblib
import uvicorn

app = FastAPI()

origins = ["*"]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/data/")
def api_data(request: Request):
    params = request.query_params
    url = f'http://pusat-api:8010/api/v1/getIOCs?{params}'
    print(url)
    r = requests.get(url)
    dictr = r.json()
    recs = dictr['data']
    # with open('saved_dictionary.pkl', 'rb') as f:
    #     recs = pickle.load(f)
    df = json_normalize(recs)
    print(df)
   
    return df_to_ml(df)




categorical_cols = ['asn', 'country_code', 'os']
le=joblib.load('/code/app/labelEncoder.joblib')
XGBoostCF = joblib.load("/code/app/XGBModel.pkl")
def df_to_ml(df):
    old_df = pd.DataFrame(df)
    df["port_count"] = 0
    df["tcp_count"] = 0
    df["udp_count"] = 0
    for index, row in df.iterrows():
        if row["port_data"] == None:
            continue
        else:
            df["port_count"][index] = len(row["port_data"])
            for i in row["port_data"]:
                if i["protocol"] == "tcp":
                    df["tcp_count"][index] += 1
                df["udp_count"][index] = df["port_count"][index] - df["tcp_count"][index]
    

    df[categorical_cols] = df[categorical_cols].apply(lambda col: le.fit_transform(col))
    X = df.drop([ "ip","port_data"], axis=1)
    y_pred = XGBoostCF.predict(X)
    print(old_df)
    old_df["status"] = y_pred
    print(old_df)
    out = old_df.to_json(orient='records', lines=True)
    
    return out
