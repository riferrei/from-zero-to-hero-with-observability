package com.riferrei.database.domain;

public class MarketEstimate {

    private  int estimate = 0;
    private  String brand;
    private  String model;
    private  int year;

    public MarketEstimate() {

    }

    public MarketEstimate(String brand, String model, int year, int estimate) {
        this.brand = brand;
        this.model = model;
        this.year = year;
        this.estimate = estimate;
    }

    public void calculateEstimate() {
        int est = (int) (Math.random() * 10000);
        setEstimate(est);
    }

    public int getEstimate() {
        return estimate;
    }

    public void setEstimate(int estimate) {
        this.estimate = estimate;
    }

    public String getBrand() {
        return brand;
    }

    public void setBrand(String brand) {
        this.brand = brand;
    }

    public String getModel() {
        return model;
    }

    public void setModel(String model) {
        this.model = model;
    }

    public int getYear() {
        return year;
    }

    public void setYear(int year) {
        this.year = year;
    }
}
