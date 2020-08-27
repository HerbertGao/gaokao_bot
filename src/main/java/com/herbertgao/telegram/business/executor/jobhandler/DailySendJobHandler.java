package com.herbertgao.telegram.business.executor.jobhandler;

import com.herbertgao.telegram.util.GaokaoBotUtil;
import com.herbertgao.telegram.database.entity.ExamDate;
import com.herbertgao.telegram.database.entity.UserTemplate;
import com.herbertgao.telegram.database.service.ExamDateService;
import com.herbertgao.telegram.database.service.UserTemplateService;
import com.xxl.job.core.biz.model.ReturnT;
import com.xxl.job.core.handler.IJobHandler;
import com.xxl.job.core.handler.annotation.XxlJob;
import com.xxl.job.core.log.XxlJobLogger;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.bots.AbsSender;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.List;

/**
 * @program: gaokao_bot
 * @description: 每日发送高考倒计时到指定频道
 * @author: HerbertGao
 * @create: 2019-06-12 16:10
 **/
@Component
public class DailySendJobHandler {

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
     *
     * @param param
     * @return
     * @throws Exception
     */
    @XxlJob(value = "dailySendJobHandler")
    public ReturnT<String> dailySendJobHandler(String param) throws Exception {

        String[] params = param.split(",");

        boolean hourly = params.length > 0 && "hourly".equals(params[0]);
        boolean force = params.length > 1 && "force".equals(params[1]);

        try {
            LocalDateTime now = LocalDateTime.now().withNano(0);

            if (hourly) {
                now = now.withMinute(0).withSecond(0).withNano(0);
            }

            List<ExamDate> examList = examDateService.getExamList(now, true);

            for (ExamDate exam : examList) {

                boolean send = force;
                long days = ChronoUnit.DAYS.between(now, exam.getExamBeginDate());
                long hours = ChronoUnit.HOURS.between(now, exam.getExamBeginDate());

                if (days > 1) {
                    if (now.getHour() == 9) {
                        send = true;
                    }
                } else if (hours <= 24 && now.isBefore(exam.getExamBeginDate())) {
                    send = true;
                }

                XxlJobLogger.log("发送倒计时: {}", send);
                if (send) {
                    UserTemplate template = userTemplateService.getDefaultTemplate();
                    String msg = GaokaoBotUtil.getCountDownString(exam, now, template.getTemplateContent());
                    if (StringUtils.isNotBlank(msg)) {
                        absSender.execute(new SendMessage(channelId, msg));
                    }
                }

            }

            return IJobHandler.SUCCESS;

        } catch (RuntimeException re) {
            return IJobHandler.FAIL;
        }

    }

}
