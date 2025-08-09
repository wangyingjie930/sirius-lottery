package com.sirius.lottery.java.application.service;

import com.sirius.lottery.java.application.dto.DrawRequest;
import com.sirius.lottery.java.application.dto.DrawResponse;
import com.sirius.lottery.java.application.dto.LotteryInstanceResponse;
import com.sirius.lottery.java.application.dto.StockActionRequest;

public interface LotteryService {

    /**
     * 核心抽奖接口
     *
     * @param req a {@link DrawRequest} object.
     * @return a {@link DrawResponse} object.
     */
    DrawResponse draw(DrawRequest req);

    /**
     * 获取活动详情，用于前端渲染
     *
     * @param instanceId a {@link java.lang.String} object.
     * @return a {@link LotteryInstanceResponse} object.
     */
    LotteryInstanceResponse getLotteryInstance(String instanceId);

    /**
     * 扣减库存
     *
     * @param req a {@link StockActionRequest} object.
     */
    void deductStock(StockActionRequest req);

    /**
     * 增加库存
     *
     * @param req a {@link StockActionRequest} object.
     */
    void increaseStock(StockActionRequest req);
}
