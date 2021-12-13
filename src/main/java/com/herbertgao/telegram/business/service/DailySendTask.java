package com.herbertgao.telegram.business.service;

import cn.hutool.core.date.LocalDateTimeUtil;
import com.herbertgao.telegram.database.entity.ExamDate;
import com.herbertgao.telegram.database.entity.UserTemplate;
import com.herbertgao.telegram.database.service.ExamDateService;
import com.herbertgao.telegram.database.service.UserTemplateService;
import com.herbertgao.telegram.util.GaokaoBotUtil;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.bots.AbsSender;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.List;

/**
 * 每日发送高考倒计时到指定频道
 *
 * @author HerbertGao
 * @date 2019-06-12
 */
@Slf4j
@Component
public class DailySendTask {

    @Autowired
    private AbsSender absSender;
    @Autowired
    private ExamDateService examDateService;
    @Autowired
    private UserTemplateService userTemplateService;
    @Value("${telegram.channel.id}")
    private String channelId;

    /**
     * 日常发送Job
     */
    @Scheduled(cron = "0 0 * * * ?")
    public void dailySendJobHandler() throws Exception {
        LocalDateTime now = LocalDateTime.now().withMinute(0).withSecond(0).withNano(0);
        log.info("到达发送时间: {}", LocalDateTimeUtil.formatNormal(now));

        List<ExamDate> examList = examDateService.getExamList(now, true);

        for (ExamDate exam : examList) {

            boolean send = false;
            long days = ChronoUnit.DAYS.between(now, exam.getExamBeginDate());
            long hours = ChronoUnit.HOURS.between(now, exam.getExamBeginDate());

            if (days > 1) {
                if (now.getHour() == 9) {
                    send = true;
                }
            } else if (hours <= 24 && now.isBefore(exam.getExamBeginDate())) {
                send = true;
            }

            log.info("是否发送倒计时: {}", send);
            if (send) {
                UserTemplate template = userTemplateService.getDefaultTemplate();
                String msg = GaokaoBotUtil.getCountDownString(exam, now, template.getTemplateContent());
                if (StringUtils.isNotBlank(msg)) {
                    absSender.execute(new SendMessage(channelId, msg));
                    log.info("倒计时发送频道: {}，倒计时信息: {}", channelId, msg);
                }
            }
        }

    }

}
